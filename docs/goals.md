# Goals for Reporting tool:

- Support most stargate enabled cosmos networks for reporting on validator income
- Support the needs of Pylon Validation/(Strangelove Ventures) tax reporting for stargate only validator income for 2021

Supported Date Range Jan 1 2021 - Now (or end at Dec 31 2021). Also should be able to pull a specific date range

### Cosmos
[Cosmos Network Info](https://github.com/cosmos/mainnet)
- [ ] `cosmoshub-3`
    - [ ] Dates
        - network start block: `1` - 12/11/2019, 16:11:34 UTC
        - pylon validator start block: `1`
        - network end block: `5200790` - 2/18/2021, 5:59:57 UTC
    - [ ] midnight blocks for dates
    - [ ] valdiator address
        - `cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf`
        - `D540AB022088612AC74B287D076DBFBC4A377A2E`
    - [ ] [archive data](https://archive.interchain.io/) (and process for spinning up)
        - [ ] TODO: spin-up archive node & query to verify dates/blocks shown above 
- [ ] `cosmoshub-4`
    - [ ] Dates
        - network start block: `5200791` - 2/18/2021, 6:00:04 UTC
        - pylon validator start block: `5200791`
    - [ ] midnight blocks for dates
    - [ ] validator address 
        - `cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf`
        - `D540AB022088612AC74B287D076DBFBC4A377A2E`
    - [ ] node: https://hub.technofractal.com:443
    - [ ] this reporting tool should work here

### Kava
[Kava network info](https://github.com/Kava-Labs/launch)
- [ ] `kava-4`
    - [ ] Dates
        - network start block: `1` - 10/15/2020, 14:00:00 UTC
        - pylon validator start block: `1`
        - network end block: `1267329` - 3/4/2021, 12:59:57 UTC
    - [ ] midnight blocks
    - [ ] validator address
        - `kavavaloper16lnfpgn6llvn4fstg5nfrljj6aaxyee9z59jqd`
        - `FF8EEC37F5F911B89FF12CB00A432157FD6E7847`
    - [ ] custom reporting tooling for this (same thing that works for `akashnet-1` should work here too)
    - [ ] archive data http://3.218.154.47:26657
- [ ] `kava-7`
    - [ ] Dates
        - network start block: `1` - 4/8/2021, 15:00:00 UTC
        - pylon validator start block: `1`
    - [ ] midnight blocks for dates
    - [ ] validator address
        - `kavavaloper16lnfpgn6llvn4fstg5nfrljj6aaxyee9z59jqd`
        - `FF8EEC37F5F911B89FF12CB00A432157FD6E7847`
    - [ ] custom reporting tooling for this (same thing that works for `akashnet-1` should work here too)
    - [ ] node: https://rpc.data.kava.io:443

### Akash
[Akash network info](https://github.com/ovrclk/net)
- [ ] `akashnet-1`
    - [ ] Dates
        - network start block: `1` - 9/5/2020, 14:00:00 UTC
        - pylon validator start block: `1`
        - network end block: `2283024` - 3/8/2021, 14:59:41 UTC
    - [ ] midnight blocks
    - [ ] validator address (verify pylon validator bond height)
        - `akashvaloper1lhenngdge40r5thghzxqpsryn4x084m9c50tcq`
        - `58B9F2517DBF8C55573FCFA9853FFF27066D86ED`
    - [ ] archive data
        - [ ] TODO: Verify above dates/blocks with archive node query
- [ ] `akashnet-2`
    - [ ] Dates
        - network start block: `1` - 3/8/2021, 15:00:00 UTC
        - pylon validator start block: `1`
    - [ ] midnight blocks
    - [ ] validator address
        - `akashvaloper1lhenngdge40r5thghzxqpsryn4x084m9c50tcq`
        - `58B9F2517DBF8C55573FCFA9853FFF27066D86ED`
    - [ ] node: https://akash.technofractal.com:443
    - [ ] this reporting tool should work here

### Sentinel
[Sentinel network info](https://github.com/sentinel-official/networks)
- [ ] `sentinelhub-1`
    - [ ] investigate if we were even in [sentinel-1](https://secretnodes.com/sentinel/chains/sentinelhub-1/validators/285643F9F72C36041189C5B9E748D665A4897EFD)
        - Added to active set @ block `146672`
        - Added to active set @ block `672511`
            - [ ] TODO: Verify pylon validator status between above blocks      
    - [ ] validator address
        - `sentvaloper1lhenngdge40r5thghzxqpsryn4x084m9mn976k`
        - `285643F9F72C36041189C5B9E748D665A4897EFD`        
    - [ ] Dates
        - network start block: `1` - 3/27/2021, 12:00:00 UTC
        - network end block: `901799` - 5/29/2021, 13:11:00 UTC
- [ ] `sentinelhub-2`
    - [ ] Dates
        - network start block: `901801` - 5/29/2021, 14:30:00 UTC
            - [ ] TODO: Verfify network upgrade block discrepency between blocks `901799-901802`
        - pylon validator start block: `901801`
    - [ ] midnight blocks
    - [ ] validator address
        - `sentvaloper1lhenngdge40r5thghzxqpsryn4x084m9mn976k`
        - `285643F9F72C36041189C5B9E748D665A4897EFD`
    - [ ] node: https://dvpn.technofractal.com:443
    - [ ] this reporting tool should work here

### Osmosis
[Osmosis network info](https://github.com/osmosis-labs/networks)
- [ ] `osmosis-1`
    - [ ] Dates
        - network start block: `1` - 6/18/2021, 17:00:00 UTC
        - pylon validator start block - `54799`  - 6/22/2021, 21:32:52 UTC
            - [ ] TODO: Verify pylon's first block in active set. querying `https://osmosis.technofractal.com/block?height=54898` shows pylon in `signatures` object, but querying `https://osmosis.technofractal.com/validators?height=54798&per_page=120` does not show pylon in `validators` object
    - [ ] midnight blocks
    - [ ] validator address
        - `osmovaloper1r2u5q6t6w0wssrk6l66n3t2q3dw2uqny4gj2e3`
        - `138FD9AB7ABE0BAED14CA7D41D885B78052A4AA1`
    - [ ] node: https://osmosis.technofractal.com:443
    - [ ] this reporting tool should work here
- [ ] Addition osmosis trading data reporting tool, scope TBD