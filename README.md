# Secure Store

This is an experimental project for encrypting content on Balena IoT devices (such as the Raspberry Pi) to avoid your content being accessed when devices are lost or stolen. Content can include both folders/files and environment variables, such as encrypted API keys stored in your Balena Cloud.

The difficultly with implementing this functionality on devices like the Raspberry Pi, is your device must be given a password to decrypt content. If you keep that password on the device, it is no longer secure. If you have to go around and enter the password on all your devices, you remove the ability to manage devices remotely. Secure Store aims to overcome this by serving the secure key (your password to unlock the device) from a Store Server device, which the Store Client devices must be able to see and access to decrypt the content stored on them. Example use cases include:

1. Run your Store Clients on the network. Plugin your Store Server to decrypt the Store Clients, and then remove the Store Server from the network again. In this instance, as soon as power is removed from the device there is only encrypted content available to the person who took your device, and no access to the Secure Server to decrypt it.
2. Run the Store Client and Store Servers on the same network, ensuring your network is secured either by protecting the Ethernet accessability or secure WiFi passwords. A compromised device when taken off your network will no longer be able to access the Store Server, and the data will remain encrypted. For a user to decrypt the data, they would need to extract the WiFi password from the device and the keys used to communicate with the server, go back to your network, connect to your network and then request the decryption key from your server using the extracted keys, significantly reducing the likelihood of a compromise.
3. Run your Store Server on a remote server online and apply an IP blocklist/whitelist to the host. Your online server will only accept requests from certain IP addresses (for example the IP address of your public network), which means devices will only decrypt data if they are requesting the key from inside your own network.

If you identify other use cases please do let us know so we can explore further iterations.

## How it works

1. You generate a Secure Store key and encrypt your content locally on your system.
2. You generate MTLS certificates for secure communication between the Server and Client. [MTLS is similar to TLS](https://www.cloudflare.com/en-gb/learning/access-management/what-is-mutual-tls/), except mTLS ensures that both parties at each end of a network connection are who they claim to be by verifying that they both have the correct private key.
3. Your encrypted content and MTLS keys are added to your Docker container
4. Your generated Secure Store key for decrypting the content is added as an environment variable on the Secure Store Server device.
5. When your Secure Store Client starts, it will look for the Secure Store Server and keep searching until it finds it.
6. Once your Client finds the Server, it will verify the MTLS certificates on each device to ensure your devices are legit and the Server will provide the Secure Store Client with the password required for decryption.
7. Your Secure Store Client will decrypt the content on the device, and decrypt any encrypted variables stored in your Balena Cloud Dashboard and then execute any application or process you choose.

The encryption only applies to environment variables and the files and folders you specify, and not the whole operating system, which avoids the performance impact that can come from full disk encryption.

# Security and encryption protocols

MTLS certificates use SHA256WithRSA and encrypted environment variables using AES-GCM with a 32bit key for AES-256. The code for generating both these components is available in this repository for your own auditing and development.

The files and folders are encrypted using [Rclone Crypt](https://github.com/rclone/rclone), which is also open source and available for your own auditing. Under the hood are [different encryption methods](https://rclone.org/crypt/) which utilise among others NaCl SecretBox based on XSalsa20 cipher and Poly1305 for integrity. It's content is encrypted using a randomly generated 1024 bit password which is unique for each device, and stored inside a configuration file. That configuration file is then encrypted with NaCl SecretBox using your own provided key generated from Secure Store and served by Secure Store Server to all your devices.

As with any encryption solution, this is not bullet proof. This project has been developed as a proof of concept designed to significantly increase the level of security of content, but does not make any guarantees.

## Demo

This repo is setup as a demo of Secure Store and can be deployed using the Deploy to Balena button below.

Instead of using two separate devices, it uses two containers merely to demonstrate the decryption process in action. A two container setup as demonstrated here would provide no security benefits as the keys are all stored on the same device, it is purely for demonstration processes. To build between devices as it is designed for, see the Installation section below.

[![balena deploy button](https://www.balena.io/deploy.svg)](https://dashboard.balena-cloud.com/deploy?repoUrl=https://github.com/maggie0002/secure-store)

When your containers start, you will see the Client looking for the Server container. When it finds it, it will decrypt the demo content and environment variables using the key it has fetched from the Server. You will then be able to go in to the container and see both the encrypted and decrypted content, as well as see a basic shell script in your terminal/logs that is started after decryption and is accessing the decrypted content.

## Setting up your own

There are many different workflows for setting up and running this project for yourself. MTLS keys and user keys should ALWAYS be stored in GitHub Secrets or other secure means, and not in GitHub repos. For the purposes of demonstrating the setup however, we will keep them in the repo for ease of understanding.

1. Start by cloning this GitHub project to your system
2. Delete all the contents of ./encrypted and ./keys
3. Replace the contents of ./source with the files you want to store encrypted
4. Generate your MTLS keys by running the command below. You will need to replace `secure-store-server` with the hostname of the device as it will be seen by the Client. For example, `secure-store.local` or `https://mystore.com`.

```
docker run -v ${PWD}/keys:/app/keys enc -generate-keys -certificate-path keys/cert.pem -key-path keys/key.pem -hostname secure-store-server
```

5. Generate a key you will use for decrypting your devices. It will be printed in your terminal.

```
docker run enc -new-key
```

6. Encrypt your files using your new key. My key is `62bca5ca7308305035f3f3a3ee270c026474ba430f79b17128021944d548807b` and is included below, replace it with your own.

```
docker run -it \
--device /dev/fuse \
--cap-add SYS_ADMIN \
-v ${PWD}/source:/app/source \
-v ${PWD}/encrypted:/app/storage \
-v ${PWD}/keys:/app/keys \
enc \
-encrypt-content ./source/. \
-config-path ./keys/encrypt.conf \
-key 62bca5ca7308305035f3f3a3ee270c026474ba430f79b17128021944d548807b
```

7. The files in ./keys should be stored in your GitHub Secrets and written to your container on build, rather than kept in your GitHub repo, but for now we will continue as is to be more transparent on how it works. The Dockerfiles included in this project provide examples of how to copy these folders in to the container in the correct places, and include the setup of Secure Store as a reference point.

8. Encrypt your environment variable using the below command, where `-string` is the variable value to encrypt, and `-key` is the key we generated earlier in step 5.

```
docker run enc -key 62bca5ca7308305035f3f3a3ee270c026474ba430f79b17128021944d548807b -string this-is-my-test-api-key
```

9. Secure Store will print an encrypted version of your passed string. We include these in our Balena Dashboard as environment variables to be decrypted later. Add them in to your dashboard in the following format:

```
ENCRYPTED_TESTVAR=d72afca84fb28581a2d9d533835a41cacb4beef534797825da7d8b80889593540077b2a4d9abb6db260108541903b6ce44a1a9
```

`ENCRYPTED_` tells Secure Store to handle these environment variables, and ignore the rest. Once the environment variable has been decrypted this prefix will be stripped and your environment variable in your container will become: `TESTVAR=this-is-my-test-api-key`

10. Finally, add your generated key as an environment variable for the Secure Store Server (NOT the client) as `STORE_PASSWORD`. For example: `STORE_PASSWORD=62bca5ca7308305035f3f3a3ee270c026474ba430f79b17128021944d548807b`. Below is an example of the Balena Dashboard setup:

You are now ready to start your devices. Enjoy!
