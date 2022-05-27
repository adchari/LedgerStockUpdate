# LedgerStockUpdate

This repository has been archived, and the project has been ported to Rust, please refer to the [SourceHut page](https://sr.ht/~adchari/ledgerdb)

This application locates any stocks you have in your [ledger-cli](https://ledger-cli.org) file, then generates a price database of those stocks compatible with the application.

### Usage

Build the go file, and run as follows:

```bash
./[name of executable] -f=[ledger file] -p=[price database file (to create or update)] -a=[Alpha Vantage API token] -b=[Name of ledger binary]
```

This should spit out a price database file, which can then be used to calculate the market value in ledger as follows:

```bash
ledger -f [ledger file] --price-db [price database file] -V bal
```
