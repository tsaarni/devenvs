// Test how mmap works when filesystem blocks
//
// Compile:
//   gcc -std=c99 -Wall -Wextra -o mmapper mmapper.c

#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/mman.h>
#include <time.h>
#include <unistd.h>

#define MAP_NUM_PAGES 100

#define INFO(fmt, ...) LOGGER(fmt, "INFO", ##__VA_ARGS__)
#define LOGGER(fmt, level, ...)                                                \
  do {                                                                         \
    time_t t = time(NULL);                                                     \
    struct tm tm = *localtime(&t);                                             \
    fprintf(stderr, "%d-%02d-%02dT%02d:%02d:%02d %s:%d " level ": " fmt "\n",  \
            tm.tm_year + 1900, tm.tm_mon + 1, tm.tm_mday, tm.tm_hour,          \
            tm.tm_min, tm.tm_sec, __FILE__, __LINE__, ##__VA_ARGS__);          \
  } while (0)

int main(int argc, char *argv[]) {
  if (argc < 2) {
    fprintf(stderr, "Usage: %s <filepath>\n", argv[0]);
    exit(EXIT_FAILURE);
  }

  const char *filepath = argv[1];
  int fd = open(filepath, O_RDWR | O_CREAT, 0666);
  if (fd == -1) {
    perror("Error opening file for writing");
    exit(EXIT_FAILURE);
  }
  INFO("File opened: %s", filepath);

  size_t page_size = sysconf(_SC_PAGESIZE);
  size_t file_size = MAP_NUM_PAGES * page_size;
  /*
    if (lseek(fd, file_size - 1, SEEK_SET) == -1) {
      perror("Error calling lseek() to 'stretch' the file");
      close(fd);
      exit(EXIT_FAILURE);
    }
    INFO("File stretched to %zu bytes", file_size);
  */
  if (write(fd, "", 1) == -1) {
    perror("Error writing last byte of the file");
    close(fd);
    exit(EXIT_FAILURE);
  }
  INFO("Last byte written to file");

  char *map =
      (char *)mmap(0, file_size, PROT_READ | PROT_WRITE, MAP_SHARED, fd, 0);
  if (map == MAP_FAILED) {
    close(fd);
    perror("Error mmapping the file");
    exit(EXIT_FAILURE);
  }
  INFO("File mmapped to memory");

  int counter = 1;
  while (1) {
    char buf[] = "Hello world %d!";
    INFO("Writing to memory page %d", counter);
    size_t offset = counter * page_size;
    sprintf(&map[offset], buf, counter);
    INFO("Done writing");

    INFO("Reading from memory page %d", counter);

    offset = counter * page_size;
    char *buf2 = (char *)malloc(page_size);
    memcpy(buf2, &map[offset], page_size);

    counter = (counter + 1) % MAP_NUM_PAGES;

    sleep(1);
  }

  return 0;
}
