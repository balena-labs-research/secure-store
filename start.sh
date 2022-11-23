#!/usr/bin/env sh

echo "The content is decrypted and the chosen programme is now running."

echo "Your decrypted environment variable is $TESTVAR. Variables like this can be accessed by your application" \
"(such as the one printing this message), but are not visible directly from the console, for example using printenv."

if [ -d "/app/storage" ] && [ -d "/app/decrypted" ]; then
    echo "You can access the client container and have a look around at the decrypted and encrypted files in" \
        "/app/decrypted and /app/storage."

    echo "Here is your list of encrypted files from /app/storage:"
    ls /app/storage

    echo "And here is a list of the same files decrypted and ready to use from /app/decrypted:"
    ls /app/decrypted
fi

echo "Enjoy!"

exec sleep infinity
