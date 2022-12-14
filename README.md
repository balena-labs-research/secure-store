# Secure Store

###### Secure Store is like a YubiKey for your network, providing encryption and remote decryption of files/folders and environment variables for entire fleets.

This is an experimental project for encrypting content on Balena IoT devices (such as the Raspberry Pi) to avoid your content being accessed when devices are lost or stolen. Content can include both folders/files and environment variables, such as encrypted API keys stored in your Balena Cloud.

The difficultly with implementing this functionality on devices like the Raspberry Pi, is your device must be given a password to decrypt content. If you keep that password on the device, it is no longer secure. If you have to go around and enter the password on all your devices, you remove the ability to manage devices remotely. Secure Store aims to overcome this by serving the secure key (your password to unlock the device) from a Store Server device, which the Store Client devices must be able to see and access to decrypt the content stored on them. Example use cases include:

1. Run your Secure Store Clients on your network. Plug your Secure Store Server in to the same network which will allow your Clients to access the key and decrypt the content store on them. Then remove the Secure Store Server from the network again. In this scenario, as soon as power is removed from the Client device there is only encrypted content available to the person who took your device, and no access to the Secure Server to decrypt it again.
2. Run the Store Client and Store Servers on the same network, ensuring your network is secured either by protecting the Ethernet accessability or secure WiFi passwords. A compromised device when taken off your network will no longer be able to access the Store Server and get the key required to decrypt the data. For a user to decrypt the data, they would need to extract the WiFi password from the device and the mTLS keys used to communicate with the server, go back to your network, connect to your network and then request the decryption key from your server using the extracted keys, significantly reducing the likelihood of a compromise.
3. Run your Store Server on a remote server online and apply an IP blocklist/whitelist to the host. Your online server will only accept requests from certain IP addresses (for example the IP address of your public network), which means devices will only decrypt data if they are requesting the key from inside your own network.
4. Encrypt the content [locally](#local-mode), and use something distinct about the operating environment as your password. For example, you could use a hash of your WiFi configuration to only decrypt when connected to the internet, or a hash of a certain part of a file unique to the OS. By compiling the requirements for the hash into the binary, it mitigates hackers attempts to access your content.

If you identify other use cases please do let us know so we can explore further iterations.

## Demo

This repo is set up as a demo of Secure Store and can be deployed using the Deploy to Balena button below.

Instead of using two separate devices, it uses two containers merely to demonstrate the decryption process in action. A two container setup as demonstrated here would provide no security benefits as the keys are all stored on the same device, it is purely for demonstration processes. To build between devices as it is designed for, see the `Setting up your own` section below.

[![balena deploy button](https://www.balena.io/deploy.svg)](https://dashboard.balena-cloud.com/deploy?repoUrl=https://github.com/balena-labs-research/secure-store)

When your containers start, you will see the Client looking for the Server container. When it finds it, it will decrypt the demo content and environment variables using the key it has fetched from the Server. You will then be able to go in to the container and see both the encrypted and decrypted content, as well as see a basic shell script in your terminal/logs that is started after decryption and is accessing the decrypted content.

## How it works

1. You encrypt your content locally on your system using Secure Store.
2. You generate mTLS certificates for secure communication between the Server and Client. [mTLS is similar to TLS](https://www.cloudflare.com/en-gb/learning/access-management/what-is-mutual-tls/), except mTLS ensures that both parties at each end of a network connection are who they claim to be by verifying that they both have the correct private key.
3. Your encrypted content is added to your Docker container
4. Your Secure Store password for decrypting the content and mTLS keys are added as environment variables on the Secure Store Server device via the Balena Cloud.
5. When your Secure Store Client starts, it will look for the Secure Store Server and keeps searching until it finds it.
6. Once your Client finds the Server, it will verify the mTLS certificates on each device to ensure your devices are legit and the Server will provide the Secure Store Client with the password required for decryption.
7. Your Secure Store Client will decrypt the content on the device, and decrypt any encrypted variables stored in your Balena Cloud Dashboard. It will then execute any application or process you choose.

The encryption only applies to environment variables and the files and folders you specify, not the whole operating system, which avoids the performance impact that can come from full disk encryption.

## Security and encryption protocols

mTLS certificates use SHA256WithRSA and encrypted environment variables use AES-GCM with a 32bit key for AES-256. The code for generating both these components is available in this repository for your own auditing and development.

The files and folders are encrypted using [Rclone Crypt](https://github.com/rclone/rclone), which is also open source and available for your own auditing. Under the hood are [different encryption methods](https://rclone.org/crypt/) for different parts of the encryption (e.g. filenames vs files themselves) which utilise among others NaCl SecretBox based on XSalsa20 cipher and Poly1305 for integrity. It's content is encrypted using a randomly generated 1024 bit password which is unique for each device, and stored inside a configuration file. That configuration file is then encrypted with NaCl SecretBox using your own provided key generated from Secure Store and served by Secure Store Server to all your devices.

As with any encryption solution, this is not bullet proof. This project has been developed as a proof of concept designed to significantly increase the level of security of content, but does not make any guarantees.

## Setting up your own

There are many different workflows for setting up and running this project for yourself. You should avoid copying unencrypted content into your containers and doing any sort of encryption in the Docker build steps as an attacker could potentially extract the unencrypted content from your Docker layers. mTLS keys and user keys should ALWAYS be stored in GitHub Secrets or other secure means, and not in GitHub repos. For the purposes of demonstrating the setup however, we will keep the keys in the repo for ease of understanding.

:grey_exclamation: Tip: If you are only looking to use the environment variable decryption (without any data storage) pass the `-env-only` flag and you can skip steps 2, 3 and 5.

1. Start by cloning this GitHub project to your system
2. Delete all the contents of `./encrypted` and `./keys`
3. Replace the contents of `./source` with the files you want to store encrypted
4. Generate your mTLS keys by running the command below. You will need to replace `secure-store-server` with the hostname of the device as it will be seen by the Client. For example, `secure-store.local` or `https://mystore.com`.

```
docker run ghcr.io/balena-labs-research/secure-store -generate-keys -hostname secure-store-server -base64
```

You will see two mTLS keys outputted to your terminal as base64 strings. Add these to your Balena Cloud account as environment variables `MTLS_CERT` and `MTLS_KEY`.

If you would rather store the mTLS keys as files on the device, remove the `-base64` flag and add `-v ${PWD}/keys:/app/keys -certificate-path keys/cert.pem -key-path keys/key.pem` to generate files for inclusion in the root directory of the Secure Store executable.

5. Encrypt your files using a [secure password](https://www.lastpass.com/features/password-generator). My example password is `my-password-eQ4al9jgPxlWDwxL6uiGdznhhVJzaVQPnkNRjwvwoTvqWpeBJJJZ` and is included below, replace it with your own.

```
docker run \
--device /dev/fuse \
--cap-add SYS_ADMIN \
-v ${PWD}/source:/app/source \
-v ${PWD}/encrypted:/app/storage \
-v ${PWD}/keys:/app/keys \
ghcr.io/balena-labs-research/secure-store \
-encrypt-content ./source/. \
-config-path ./keys/encrypt.conf \
-password my-password-eQ4al9jgPxlWDwxL6uiGdznhhVJzaVQPnkNRjwvwoTvqWpeBJJJZ
```

The files generated in `./keys` should be stored in your GitHub Secrets and written to your container on build rather than kept in your GitHub repo but for now we will continue by copying them in to the containers through the Dockerfiles to be more transparent on how it works.

6. Encrypt your environment variable using the below command, where `-string` is the variable value to encrypt, and `-password` is the same password as chosen above for encrypting files.

```
docker run ghcr.io/balena-labs-research/secure-store -password my-password-eQ4al9jgPxlWDwxL6uiGdznhhVJzaVQPnkNRjwvwoTvqWpeBJJJZ -string this-is-my-new-test-api-key
```

Secure Store will print an encrypted version of your passed string. We include these in our Balena Dashboard as environment variables to be decrypted later. Add them in to your dashboard in the following format:

```
ENCRYPTED_TESTVAR=d72afca84fb28581a2d9d533835a41cacb4beef534797825da7d8b80889593540077b2a4d9abb6db260108541903b6ce44a1a9
```

The `ENCRYPTED_` prefix tells Secure Store to handle these environment variables. Once the environment variable has been decrypted this prefix will be stripped and your environment variable in your container will become: `TESTVAR=this-is-my-test-api-key`

7. Finally, add your password as an environment variable for the Secure Store Server (NOT the client) as `STORE_PASSWORD`. For example: `STORE_PASSWORD=my-password-eQ4al9jgPxlWDwxL6uiGdznhhVJzaVQPnkNRjwvwoTvqWpeBJJJZ`. Below is an example of a Balena Cloud Dashboard setup:

<img width="1188" alt="env-examples" src="https://user-images.githubusercontent.com/64841595/199483420-6dc7b50e-b8be-4503-9348-95459cee0b80.png">

You are now ready to start your devices. Enjoy!

## Local mode

You can start the client in local mode by passing `-local your-password`. This will create the mount and not rely on the server. Your password could however, be extracted from your device if you keep the two together without additional security measures. You should consider compiling the password used, or the means of generating the password into the Golang binary. 

## Using the encrypted storage for volumes

When the device decrypts the content, it will create a folder called `decrypted` on the path specified (defaults to `./`). In this folder you can write new content and it will be included in your encrypted storage. The `decrypted` folder is a mount, that points to the encrypted content behind it stored in a folder called `storage`. If you make `storage` a volume, then the encrypted content will be persistent.
