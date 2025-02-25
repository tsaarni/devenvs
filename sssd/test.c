// Compile with: gcc test.c -o test
#include <stdio.h>
#include <libintl.h>
#include <locale.h>

int main() {
    if (setlocale(LC_ALL, "") == NULL) {
        perror("setlocale");
        return 1;
    }

    if (textdomain("sshd") == NULL) {
        perror("textdomain");
        return 1;
    }

    if (bindtextdomain("sssd", "/usr/share/locale") == NULL) {
        perror("bindtextdomain");
        return 1;
    }
    char* message = "Your password has expired. You have %1$d grace login(s) remaining.";
    char* translated_message = dgettext("sssd", message);
    printf("%s\n", translated_message);

    return 0;
}
