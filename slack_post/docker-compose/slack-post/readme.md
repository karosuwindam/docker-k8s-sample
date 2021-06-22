要件定義

github.com/ashwanthkumar/slack-go-webhookを使用して作成する

Client-goを使用して、pod情報から自動で、監視Podを検出する

Annotations内のキー名slackpostがtrueの値のものを検出したら、そのPODは監視対象とする

監視対象となった機器は、http://IP名:ポート/jsonの値を取得する

jsonの形式は以下の通り、

```
[
    {
        "Sennser":"名称",
        "Type":"種類",
        "Data":"値"
    }
]
```