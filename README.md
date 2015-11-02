# libsecret Docker Volume Driver Plugin
This is a Docker Volume Driver plugin for [libsecret](https://github.com/ehazlett/libsecret).

# Usage
To start the plugin, run the following:

This example uses [Vault](https://vaultproject.io)

```bash
sudo docker-volume-libsecret --addr <vault-address> --backend vault --store-opt token=<vault-token>
```

Replace `<vault-address>` and `<vault-token>` with the address and token
to your Vault instance.

Then, create a secret.  For example, set the key `secret/app/db/username` to `root`.

You can then run a container mounting any part of the path and the container
will be able to read the secret.

```bash
docker run -ti --rm --volume-driver libsecret -v secret/app:/app alpine ash

/ # cat /app/db/username
root
/ #
```

Note: You will not be able to browse (`ls`) the mounted directory (`/app`).  
This is by design and for security to prevent unauthorized viewing.  The 
application must know the path to the secrets.
