# Home Controller
[Nature Remo](https://nature.global/nature-remo/)から信号を送るコマンドラインツールです。

send_signalモードとtemp_controlモードの2つのモードがあります。
- send_signalモード: 選択肢からインタラクティブに操作するモード
- temp_controlモード: 室温が規定の範囲内になるようにエアコンの設定温度を自動変更するモード

## 使い方
### 事前準備
Nature Remoを設置し、各家電のリモコンを設定します。
https://shop.nature.global/

Nature RemoのAPIキーを取得します。
https://home.nature.global/

Nature RemoのAPIキーを環境変数に設定します。
```bash
export NATURE_REMO_SECRET="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

### send_signalモード
```bash
go run main.go -nature_api_secret=${NATURE_REMO_SECRET} -action_mode="send_signal"
```

### temp_controlモード
```bash
go run main.go -nature_api_secret=${NATURE_REMO_SECRET} -action_mode="temp_control" -device_name="Remo"
```