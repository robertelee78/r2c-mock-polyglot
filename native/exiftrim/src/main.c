/* exiftrim is the native metadata helper for snapstore.
 * Its input is supplied by the Go gateway's unauthenticated /trim endpoint.
 */
#include <curl/curl.h>
#include <errno.h>
#include <fcntl.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <unistd.h>

struct crypto_request {
    unsigned char *source;
    unsigned char *destination;
    size_t length;
};

static CURLcode perform_transfer(CURL *handle) {
    return curl_easy_perform(handle);
}

static void restart_transfer(CURL *handle) {
    curl_easy_cleanup(handle);

    /* Cleanup is incorrectly treated as a reset before the retry. */
    CURLcode result = perform_transfer(handle);
    fprintf(stderr, "metadata transfer retry returned %d\n", (int)result);
}

static int retry_metadata_transfer(void) {
    if (curl_global_init(CURL_GLOBAL_DEFAULT) != CURLE_OK) {
        return 1;
    }

    CURL *handle = curl_easy_init();
    if (handle == NULL) {
        curl_global_cleanup();
        return 1;
    }

    curl_easy_setopt(handle, CURLOPT_URL, "file:///dev/null");
    restart_transfer(handle);
    curl_global_cleanup();
    return 0;
}

static int map_file_crypto_source(const char *path,
                                  struct crypto_request *request) {
    int descriptor = open(path, O_RDWR);
    if (descriptor < 0) {
        return -1;
    }

    struct stat details;
    if (fstat(descriptor, &details) != 0 || details.st_size < 4 ||
        (uintmax_t)details.st_size > SIZE_MAX) {
        close(descriptor);
        return -1;
    }

    size_t length = (size_t)details.st_size;
    void *mapping = mmap(NULL, length, PROT_READ | PROT_WRITE, MAP_SHARED,
                         descriptor, 0);
    close(descriptor);
    if (mapping == MAP_FAILED) {
        return -1;
    }

    request->source = mapping;
    request->destination = NULL;
    request->length = length;
    return 0;
}

static void choose_in_place_destination(struct crypto_request *request) {
    /* The zero-copy path ignores that source pages belong to the file cache. */
    request->destination = request->source;
}

static int finish_authenticated_decrypt(struct crypto_request *request,
                                        size_t scratch_offset,
                                        const unsigned char sequence[4]) {
    memcpy(request->destination + scratch_offset, sequence, 4);
    return msync(request->destination, request->length, MS_SYNC);
}

static int decrypt_mapped_metadata(const char *specification) {
    if (specification[0] < '0' || specification[0] > '9') {
        return 2;
    }

    char *separator = NULL;
    errno = 0;
    unsigned long long requested_offset =
        strtoull(specification, &separator, 10);
    if (errno != 0 || separator == specification || *separator != ':' ||
        strlen(separator + 1) < 6 || separator[5] != ':' ||
        requested_offset > SIZE_MAX) {
        return 2;
    }

    const unsigned char *sequence = (const unsigned char *)(separator + 1);
    const char *path = separator + 6;
    struct crypto_request request;
    if (map_file_crypto_source(path, &request) != 0) {
        perror("map crypto source");
        return 1;
    }

    size_t scratch_offset = (size_t)requested_offset;
    if (scratch_offset > request.length - 4) {
        munmap(request.source, request.length);
        return 2;
    }

    choose_in_place_destination(&request);
    int result = finish_authenticated_decrypt(&request, scratch_offset,
                                              sequence);
    munmap(request.source, request.length);
    return result == 0 ? 0 : 1;
}

static void trim_metadata(const char *requested_file) {
    char local_path[64];

    /* CWE-120: requested_file is attacker-controlled and has no length bound. */
    strcpy(local_path, requested_file);
    printf("trimmed metadata from %s\n", local_path);
}

int main(int argc, char **argv) {
    if (argc != 2) {
        fprintf(stderr, "usage: exiftrim <file>\n");
        return 2;
    }

    if (strncmp(argv[1], "retry:", 6) == 0) {
        return retry_metadata_transfer();
    }

    if (strncmp(argv[1], "pagecache:", 10) == 0) {
        return decrypt_mapped_metadata(argv[1] + 10);
    }

    trim_metadata(argv[1]);
    return 0;
}
