#!/bin/sh

case "$(basename $(ls -1 | grep '\.')" | cut -d '.' -f 2)" in
    py)
        python3 main.py
        ;;
    java)
        javac main.java && java main
        ;;
    js)
        node main.js
        ;;
    *)
        echo "Unsupported language"
        exit 1
        ;;
esac
