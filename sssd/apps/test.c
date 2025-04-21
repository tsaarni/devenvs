// Compile with: gcc test.c -o test
#include <stdio.h>
#include <libintl.h>
#include <locale.h>

int main() {
	char *locale = setlocale(LC_ALL, "");
	if (locale == NULL) {
		perror("setlocale");
		return 1;
	}
	printf("locale set to %s\n", locale);

	if (textdomain("sshd") == NULL) {
		perror("textdomain");
		return 1;
	}

	if (bindtextdomain("sssd", "/usr/share/locale") == NULL) {
		perror("bindtextdomain");
		return 1;
	}
	char *message =
		"Your password has expired. You have %1$d grace login(s) remaining.";
	char* translated_message = dgettext("sssd", message);
    printf("%s\n", translated_message);

    return 0;
}
