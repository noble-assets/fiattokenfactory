cd proto
buf generate --template buf.gen.gogo.yaml
buf generate --template buf.gen.pulsar.yaml
cd ..

cp -r github.com/circlefin/noble-fiattokenfactory/* ./
rm -rf github.com
rm -rf circle
