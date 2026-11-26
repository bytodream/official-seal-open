#!/usr/bin/env bash

cd "$(git rev-parse --show-toplevel)" && \
[ -x ./build-out/dist/official_seal_darwin_arm64 ] || ./script/cmd-build.sh && \
./build-out/dist/official_seal_darwin_arm64 \
  -attr-text '{"sealCategory":"enterprise","sealSpec":"CIRCULAR_42","content":"北京示例信息技术有限公司","foot":"合同专用章","enterpriseCode":"110108123456789"}' \
  -name my-company \
  -out ./build-out/my-seals \
  -seed 42 \
  -license 'AgEAy7qZ7+eVbLwP0H7HR0MUd9IppoSqwvSgdr6IWTD8yzN5gesf5PgNMi4zkX58KTJB9Kk5F0YFtv7Pfr1TmoaDwv3ksIDMU63RNYEt8HAuM4j4r2b0kC52E96SgSmamG1aw3Gw2lDKDT+wdkzwl71ZuOx2b7prCcecTELdv9nJ+tyDJqYj79p1R1O0CzAHL1Ra1jmhzj7RVKwJNg5Ua7+cIzeyFp+pGfkTAuO0cHxoiT+ag5O+bnIgZdzSNMCAlLdrJt0s5qipt99dGzAjRTs4cDOBozMRE5YcFfG09IRjaMM4dpmSVTECHUonZ2RvyYpjNiLM2IsZo0Hzt1V26qUM6+97SwrB2tRnN9PveYzWZo5MKQRfyCWgYVQ+ORx8i48rg4iBbd0s0vjCwstLR83zeaE4/hnZL4uZzkdK3THT2lUSghfEYCdeGvXDkQF7rTsiNmUoJGEkU+j8wdo1c2f1I0zZ8fFFSYD0t6RbHlgagyRoHQF6145sOizJSgTXG68RP3EbabACcJ4FXglD7S/YFniDCg9tbwEchLUZyALALdERVyY7y8gbVYlfowJdeTvBA0Zvjp6SvrdjQK+5dPcncFTRN6vboghAcExbKDKYEWzn30Kt6lE7xMvxmjnnES5/LW/Py30mtGD8ULgFSmF0MWuTTTvTWuhW3awHfpA5Ku4pxmq2llTxjpd8EaAZyLfCTGbg' && \
ls -lh ./build-out/my-seals