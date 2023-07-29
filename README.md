meをチェック
curl -X GET "https://api.nature.global/1/users/me" -H "accept: application/json" -H "authorization: Bearer ${NATURE_REMO_SECRET}"

nature remo デバイスを一覧
curl -X GET "https://api.nature.global/1/devices" -H "accept: application/json" -H "authorization: Bearer ${NATURE_REMO_SECRET}" | python -m json.tool --no-ensure-ascii

登録電化製品一覧
curl -X GET "https://api.nature.global/1/appliances" -H "accept: application/json" -H "authorization: Bearer ${NATURE_REMO_SECRET}" | python -m json.tool --no-ensure-ascii
