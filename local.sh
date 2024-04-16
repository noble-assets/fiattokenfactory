alias ftfd=./simapp/build/simd

for arg in "$@"
do
    case $arg in
        -r|--reset)
        rm -rf .ftf
        shift
        ;;
    esac
done

if ! [ -f .forwarding/data/priv_validator_state.json ]; then
  ftfd init validator --chain-id "ftf-1" --home .ftf &> /dev/null

  ftfd keys add validator --home .ftf --keyring-backend test &> /dev/null
  ftfd genesis add-genesis-account validator 1000000ustake --home .ftf --keyring-backend test

  TEMP=.ftf/genesis.json
  touch $TEMP && jq '.app_state.bank.denom_metadata = [{"description":"USD Coin","denom_units":[{"denom":"uusdc","exponent":0,"aliases":["microusdc"]},{"denom":"usdc","exponent":6,"aliases":[]}],"base":"uusdc","display":"usdc","name":"usdc","symbol":"USDC"}]' .ftf/config/genesis.json > $TEMP && mv $TEMP .ftf/config/genesis.json
  touch $TEMP && jq '.app_state.fiattokenfactory.mintingDenom = {"denom":"uusdc"}' .ftf/config/genesis.json > $TEMP && mv $TEMP .ftf/config/genesis.json
  touch $TEMP && jq '.app_state.staking.params.bond_denom = "ustake"' .ftf/config/genesis.json > $TEMP && mv $TEMP .ftf/config/genesis.json

  ftfd genesis gentx validator 1000000ustake --chain-id "ftf-1" --home .ftf --keyring-backend test &> /dev/null
  ftfd genesis collect-gentxs --home .ftf &> /dev/null

  sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' .ftf/config/config.toml
fi

ftfd start --home .ftf
