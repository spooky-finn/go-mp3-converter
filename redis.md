redis result with error

```json
{
  "progress": 0,
  "originalUrl": "https://www.youtube.com/watch?v=Yy4pcKn0Y_k",
  "taskName": "ae0800368a5e09412f04e0bc298851a19a4060cc",
  "status": "error",
  "origin": "it.onlinevideoconverter.pro",
  "pushedAt": 1685035783,
  "instance": "1a53473e-9aa0-4cd1-b681-52190c988b5b",
  "meta": {
    "title": "The Yussef Dayes Experience - Live At Joshua Tree (Presented by Soulection)",
    "source": "https://www.youtube.com/watch?v=Yy4pcKn0Y_k",
    "duration": "18:19",
    "tags": "yussefdayes,roccopalladino,charliestacey,welcometothehills,jazz,funk,hiphop,Yussef Dayes,Brownswood,Soulection,Joshua Tree,Indie Soul,British Jazz,Contemporary Jazz,Alexander Bourt,Drums,Bass,Percussion,Gilles Peterson,Saxophone,Synthesizer,Keys,UK contemporary Jazz,UK Jazz,Desert,Golden Hour,Live,Live session",
    "subtitle": {
      "token": "f261cbe438a2a665cc1b1e1bb6cf7161",
      "language": ["en"]
    }
  },
  "thumb": "https://i.ytimg.com/vi/Yy4pcKn0Y_k/hqdefault.jpg",
  "startAt": 1685035783.533,
  "errorMsg": "Duration limited! 1099/600"
}
```

redis good result

```json
{
  "progress": 100,
  "originalUrl": "https://www.youtube.com/watch?v=bWkP-mt_bOk",
  "taskName": "f5e65fd6d6347b5731f5288ab05eca311cef3d64",
  "status": "ready",
  "origin": "en1.onlinevideoconverter.pro",
  "pushedAt": 1685036156,
  "instance": "e92d0b36-7a16-479b-894c-783d9051b5e3",
  "meta": {
    "title": "Freemasons feat. Sophie Ellis-Bextor - Heartbreak (Make Me A Dancer)",
    "source": "https://www.youtube.com/watch?v=bWkP-mt_bOk",
    "duration": "3:35",
    "tags": "Freemasons,Sophie,Ellis-Bextor,Heartbreak,Loaded,Records,Make Me A Dancer",
    "subtitle": {
      "token": "01b1b0adc3ecf45de733b6b58d31766b",
      "language": ["en"]
    }
  },
  "thumb": "https://i.ytimg.com/vi/bWkP-mt_bOk/hqdefault.jpg",
  "fileSize": 3351223,
  "startAt": 1685036157.085,
  "duration": 215,
  "filename": "f5e65fd6d6347b5731f5288ab05eca311cef3d64.mp3",
  "stopAt": 1685036161.966
}
```

new task for converter
```json
{
  "taskName": "8bb75a0b05131a0037801f742f8d4243acf99375",
  "url": "https://youtu.be/tJqYM8b3l5k",
  "originalUrl": "https://youtu.be/tJqYM8b3l5k",
  "origin": "it.onlinevideoconverter.pro",
  "pushedAt": 1685036346
}
```