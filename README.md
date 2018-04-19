img_picker
==========

# Install

1. Install [facedetect](https://github.com/wavexx/facedetect)
2. `go build`

# Usage

`GET /data`

取得曾經上傳過的圖檔，和他的人臉位置長寬

`GET /upload`

圖檔上傳介面

`POST /upload`

上傳並分析人臉，會回傳結果
