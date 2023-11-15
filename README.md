# fxpixel [![Build and release](https://github.com/danesparza/fxpixel/actions/workflows/release.yaml/badge.svg)](https://github.com/danesparza/fxpixel/actions/workflows/release.yaml)
REST service for RGB(W) LED lighting effects on demand from Raspberry Pi. Made with ❤️ for makers, DIY craftsmen, and professional prop designers everywhere

## To install
To install on a Raspberry Pi

### Step 1
**Option 1:** Use the automated script to install the repo, the repo key and then update:
```
wget https://danesparza.github.io/package-repo/prereq.sh -O - | sh
```

**Option 2:** I don't trust you.  I'll run the commands myself:
```
curl -s --compressed "https://danesparza.github.io/package-repo/KEY.gpg" | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/package-repo.gpg >/dev/null
sudo curl -s --compressed -o /etc/apt/sources.list.d/package-repo.list "https://danesparza.github.io/package-repo/package-repo.list"
sudo apt update
```

### Step 2
Now that the repo is installed, you can install the package
```
sudo apt install fxpixel
```
