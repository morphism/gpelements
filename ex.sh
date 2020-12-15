# Some examples

# curl -c cookies 'https://www.space-track.org/basicspacedata/query/class/gp/EPOCH/%3Enow-30/NORAD_CAT_ID/270000--339999/orderby/NORAD_CAT_ID/format/json' > now.json

cat tmp/now.json |
    tletool transform -emit xml |
    tletool transform -emit kvn |
    tletool transform -emit csv |
    tletool transform -emit json |
    tletool prop -from 2020-12-15T12:00:00Z |
    tail -4 |
    jq -r -c '{"NORAD":.Norad,"LLA":.LLA}'

cat tmp/now.json |
    tletool prop -from 2020-12-15T12:00:00Z |
    tail -4 |
    jq -r -c '{"NORAD":.Norad,"LLA":.LLA}'

cat tmp/now.json |
    tletool transform -emit xml |
    tletool walk |
    tletool prop -from 2020-12-15T12:00:00Z |
    tail -4 |
    jq -r -c '{"NORAD":.Norad,"LLA":.LLA}'

cat tmp/now.json |
    tletool prop -from 2020-12-15T12:00:00Z -to 2020-12-15T12:01:00Z -interval 10s |
    tail -6 |
    jq -r -c '{"NORAD":.Norad,"LLA":.LLA}'

cat tmp/now.json |
    tletool prop -from 2020-12-15T12:00:00Z -higher-precision=false |
    tail -4 |
    jq -r -c '{"NORAD":.Norad,"LLA":.LLA}'

