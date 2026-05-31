# Home Controller

[Nature Remo](https://nature.global/nature-remo/) を使って家電を制御するコマンドラインツールです。

## モード

- `temp_control`: 室温が規定の範囲内になるようにエアコンの設定温度を自動変更する
- `notify_temp`: 室温がしきい値を大きく外れたときに `ntfy` へ通知する

## 使い方

### 事前準備

1. Nature Remo を設置し、各家電のリモコンを設定します。
   - [Nature Store](https://shop.nature.global/)
2. Nature Remo の API キーを取得します。
   - [Nature Home](https://home.nature.global/)
3. API キーを環境変数に設定します。

```bash
export NATURE_API_SECRET="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

### `temp_control` モード

室温と現在のエアコン設定を見て、設定温度を自動調整します。設定を変更したときは Slack に結果を通知し、しきい値を大きく外れた場合は補助的に `ntfy` に緊急通知します。

```bash
go run main.go \
  -nature_api_secret="${NATURE_API_SECRET}" \
  -action_mode="temp_control" \
  -device_name="Remo" \
  -tooHotThreshold=27.5 \
  -tooColdThreshold=24.5 \
  -preparationThreshold=0.5 \
  -minimumTemperatureSetting=22.0 \
  -maximumTemperatureSetting=30.0 \
  -slackToken="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
  -slackChannel="#xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
  -ntfyUrl="https://ntfy.sh/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

### `notify_temp` モード

エアコンが ON のときだけ、`temp_control` と同じ基準で、しきい値を大きく外れた室温を `ntfy` に通知します。

```bash
go run main.go \
  -nature_api_secret="${NATURE_API_SECRET}" \
  -action_mode="notify_temp" \
  -device_name="Remo" \
  -tooHotThreshold=27.5 \
  -tooColdThreshold=24.5 \
  -ntfyUrl="https://ntfy.sh/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```
