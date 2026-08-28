/* exiftrim is the native metadata helper for snapstore.
 * Its input is supplied by the Go gateway's unauthenticated /trim endpoint.
 */
#include <stdio.h>
#include <string.h>
#include <curl/curl.h>

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

    trim_metadata(argv[1]);
    return 0;
}
