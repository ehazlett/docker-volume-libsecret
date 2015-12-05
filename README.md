# libsecret Docker Volume Driver Plugin
This is a Docker Volume Driver plugin for [libsecret](https://github.com/ehazlett/libsecret).

# Usage
To start the plugin, run the following:

This example uses [Vault](https://vaultproject.io)

Create a Vault configuration file.  Note: this is just for testing. Do not use
in production.

```
backend "inmem" {}

listener "tcp" {
    address = "0.0.0.0:8200"
    tls_disable = 1

}
```

Start Vault

```bash
docker run -p 8200:8200 --name vault -ti -v /path/to/config.hcl:/vault.hcl --rm jess/vault server -config /vault.hcl
```

This will start Vault and listen on port 8200.

Next, use Docker exec to configure Vault:

```bash
docker exec -ti vault ash
/ # export VAULT_ADDR=http://127.0.0.1:8200
```

Then run `vault init` to initialize:

```bash
/ # vault init
Key 1: e51d7dab67f5f1d6ab396dba1f56cf6f29a626122c15e90ba0cddb2314d5bbe101
Key 2: ee7486a8be40a8f04c84533674b1e792440f0f5e7bddf03b3e5068ffddf8b73302
Key 3: 19340efb1645dad09adea510ed891d31e225d60a81ea0bdd107dad86d39265fe03
Key 4: d59de8adc555c4f931e1ec6638417e7c44e07b8084b5f15e681e434cce4b515c04
Key 5: 22dd60fe6d50b6d9e7bb1a40a17984dfe2caa2d47e820ab846338635c021839105
Initial Root Token: 721b2a50-629d-a47c-fe3b-95e766080087

Vault initialized with 5 keys and a key threshold of 3. Please
securely distribute the above keys. When the Vault is re-sealed,
restarted, or stopped, you must provide at least 3 of these keys
to unseal it again.

Vault does not store the master key. Without at least 3 keys,
your Vault will remain permanently sealed.
```

Use the Token and keys to unseal the vault:

```bash
/ # export VAULT_TOKEN=721b2a50-629d-a47c-fe3b-95e766080087
```

Then run `vault unseal` 3 times to unlock using the keys above:

```bash
/ # vault unseal
Key (will be hidden): 
Sealed: true
Key Shares: 5
Key Threshold: 3
Unseal Progress: 1
```

Use the `vault` command to add secrets:

```bash
/ # vault write secret/prod/redis value=foopass
Success! Data written to: secret/prod/redis

/ # vault read secret/prod/redis
Key             Value
lease_duration  2592000
value            foopass
```

Now that you have Vault, start the Volume driver plugin:

Note: due to a limitation by design, you cannot run this plugin in a container
since the fuse mounts will not be propagated.  This must run on the host.
Future updates to Docker may remove this limitation.

```bash
sudo docker-volume-libsecret --addr <vault-address> --backend vault --store-opt token=<vault-token>
```

Replace `<vault-address>` and `<vault-token>` with the address and token
to your Vault instance.

You can then run a container mounting any part of the path and the container
will be able to read the secret.

```bash
docker run -ti --rm --volume-driver libsecret -v secret/prod:/secrets alpine ash

/ # cat /secrets/redis
foopass
```

Note: You will not be able to browse (`ls`) the mounted directory (`/secrets`).  
This is by design and for security to prevent unauthorized viewing.  The 
application must know the path to the secrets.
