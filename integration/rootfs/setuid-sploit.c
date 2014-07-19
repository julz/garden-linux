#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <unistd.h>

int main(int argc, char *argv[]) {
	if (setuid(0) < 0) {
		perror("setuid");
		abort();
	}

	if (execv(argv[1], &argv[1]) < 0) {
		perror("execl");
		abort();
	}

	return 0;
}
