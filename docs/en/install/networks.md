# Networks

This document describes how to set up your own QOS network, single node or cluster mode.

## Single-node
* init

run [initialization](../command/qosd.md#initialization) command:
```bash
$ qosd init --moniker moniker --chain-id qos-test
{
 "moniker": "moniker",
 "chain_id": "qos-test",
 "node_id": "66853240dc1b26e6f6b35afcf008658823542076",
 "gentxs_dir": "",
 "app_message": {
  "accounts": null,
  "mint": {
   "params": {
    "total_amount": "10000000000",
    "total_block": "6307200"
   }
  },
  "stake": {
   "params": {
    "max_validator_cnt": 10,
    "voting_status_len": 100,
    "voting_status_least": 50,
    "survival_secs": 600
   },
   "validators": null
  },
  "qcp": {
   "ca_root_pub_key": null
  },
  "qsc": {
   "ca_root_pub_key": null
  }
 }
}
```
the configuration file will be generated by default in the `$HOME/.qosd` directory.

* add-genesis-accounts

using `qosd add-genesis-accounts` command to add accounts:

> using `qoscli keys add ` to create address information in the local keystore.

```bash

$ qoscli keys add qosInitAcc
Enter a passphrase for your key:
Repeat the passphrase:

$ qoscli keys list

NAME:   TYPE:   ADDRESS:                                                PUBKEY:
qosInitAcc      local   qosacc1lly0audg7yem8jt77x2jc6wtrh7v96hgve8fh8  qosaccpub1zcjduepquqf6k5n8y88wywmt40376m5n9zcwsz5kmnl95j7zw2l7mazf22sq3wvmur

```
visit [qoscli keys](../command/qoscli.md#keys) for more information.

according to [add genesis accounts](../command/qosd.md#genesis-accounts) for account initialization:
```bash
$ qosd add-genesis-accounts qosacc1lly0audg7yem8jt77x2jc6wtrh7v96hgve8fh8,1000000qos
```

* config-root-ca

root CA will be used in [QSC](../spec/qsc) and [QCP](../spec/qcp), no need to config if where is no related business. View [CA doc](../spec/ca.md) to learn more about CA.

using `qosd config-root-ca` command to initialize root CA public key:

```bash
$ qosd config-root-ca --qcp <qcp-root.pub> --qsc <qsc-root.pub>
```

* create-validator

using `qosd gentx` and `qosd collect-gentxs` commands to initialize the validator to the configuration file. 

```bash
$ qosd gentx --moniker validatorName --owner qosacc1lly0audg7yem8jt77x2jc6wtrh7v96hgve8fh8 --tokens 10
```

main parameters:
- `--owner`         owner address
- `--moniker`       name of this node
- `--logo`          logo
- `--website`       website
- `--details`       description
- `--tokens`        binding tokens
- `--compound`      delegation `coupound`, default `false`

for more information about `gentx` please visit [gentx doc](../command/qosd.md#gentx).

run
```bash
$ qosd collect-gentxs
```
will write the create validator transaction to the `genesis.json` file.

* start
```bash
$ qosd start --log_level debug
```
If everything is ok, you will see the console output block information.

## Cluster

* qosd testnet
[qosd-testnet](../command/qosd.md#testnet) command can batch generate multiple validators node configuration information of a test network.

assume IP of the first machine is `192.168.1.100`:
```bash
$ qosd testnet --v 4 --name capricorn --starting-ip-address 192.168.1.100
Successfully initialized 4 node directories

```
The `mytestnet` folder will be generated in the current directory, and the node 0-3 configuration files will be stored separately.
Notice that the `priv_validator_owner.json` is the private key of the validator owner, and it can be import to the local keystore by using `qoscli keys import` command.

* start
Be sure to install QOS correctly on all four machines in accordance with [Installation Instructions] (installation.md) before starting.
Copy node0-3 to a different machine and execute them separately:
```bash
$ qosd start --home <directory_for_config_and_data>
```