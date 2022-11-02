#!/usr/bin/env sh
while true
do
    echo "The content is decrypted and now the programme is running. You can access the client container \
and have a look around at the decrypted and encrypted files in /app/decrypted and /app/storage."

    echo "Your decrypted environment variable is $TESTVAR. Variables like this can be accessed by your application \
(such as this one you are reading), but are not visible directly from the console, for example using printenv."

    echo "Here is your list of encrypted files:"
    ls /app/storage

    echo "And here is your list of decrypred files:"
    ls /app/decrypted

    echo "Enjoy!"

    sleep infinity
done
