# OVERWRITE

This is a project for automatization working with a poll in telegram's chat. This project in telegram's bot do next thing:

1. Get random fact from `url`
2. Concatenate it with `question`
3. Send poll message in chat (with `chat_id`) and pin it
4. Close and unpin the poll, summarize result and send result message in chat

## LIBRARY

I use next libs for project:

1. Sheduler: `github.com/carlescere/scheduler`
2. TelegramBot: `gopkg.in/telebot.v3`

## SETTINGS AND INSTALLS

Before to build a programm you need:

0. Install [go](https://go.dev/doc/install)
1. Download the code to your folder
2. Create **config.yaml** in *your_download_folder/app/config/* directive like this:

```yaml
url: "https://randstuff.ru/fact/fav/"
path_to_pics: ""
# poll data
poll:
  question: "Завтра собираемся?"
  answersYes: [Да, Yes, Y, YES! YES! YES!, Конечно да, ➕, 👍, 🤝, 🫶, 🌞, 🗿, ✅, 🔋, 🩷, 🫰, 🆗, 💯]
  answersNo: [Нет, No, N, NO! NO! NO!, Конечно нет, ➖, 👎, ✋, 😐, 🌚, 🚧, ⛔️, 🪫, 💔, 🖕, 🛑, 🚫]
# bot data
bot_secure:
  # bot_token in base64
  bot_token: "MTIzNDVxd2VydHk=" #12345qwerty for example
  # hash256 of bot_token
  bot_hash: "b554b3ad18c9605a700f350f2e3708af9ca03857afbefde4a429ec0d6b1a9965"
  chat_id: -123456789
  # upd time for reading new message in seconds
  upd_time: 10
  # stop bot after first action
  exit_after_exec: true
```

3. Run **install.sh** in current programm folder:

```bash
./install.sh
```

### For uninstalling programm run uninstall.sh in current programm folder

```bash
./uninstall.sh
```
