# Home Controller

[Nature Remo](https://nature.global/nature-remo/) を使って家電を制御するコマンドラインツールです。

## モード

- `temp_control`: 室温が規定の範囲内になるようにエアコンの設定温度を自動変更する
- `notify_temp`: 室温が閾値を超えたときに Slack 通知する

## 使い方

### 事前準備

1. Nature Remo を設置し、各家電のリモコンを設定します。
   - [Nature Store](https://shop.nature.global/)
2. Nature Remo の API キーを取得します。
   - [Nature Home](https://home.nature.global/)
3. API キーを環境変数に設定します。

```bash
export NATURE_REMO_SECRET="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

### `temp_control` モード

```bash
go run main.go \
  -nature_api_secret="${NATURE_REMO_SECRET}" \
  -action_mode="temp_control" \
  -device_name="Remo" \
  -tooHotThreshold=27.5 \
  -tooColdThreshold=24.5 \
  -preparationThreshold=0.5 \
  -minimumTemperatureSetting=22.0 \
  -maximumTemperatureSetting=30.0 \
  -slackToken="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
  -slackChannel="#xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

### `notify_temp` モード

```bash
go run main.go \
  -nature_api_secret="${NATURE_REMO_SECRET}" \
  -action_mode="notify_temp" \
  -device_name="Remo" \
  -tooHotThreshold=27.5 \
  -tooColdThreshold=24.5 \
  -slackToken="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
  -slackChannel="#xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```
