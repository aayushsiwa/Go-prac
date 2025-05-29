# 🔤 Morse Code Converter

A modern web application built with **Go** that converts text to Morse code and vice versa. Features a clean, responsive UI with a powerful REST API backend.

![Morse Code Converter](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Status](https://img.shields.io/badge/Status-Active-brightgreen?style=for-the-badge)

## ✨ Features

-   **Bidirectional Conversion**: English ↔ Morse Code
-   **Modern Web Interface**: Responsive design with gradient backgrounds
-   **REST API**: JSON-based API for integration with other applications
-   **Comprehensive Character Support**: A-Z, 0-9, and common punctuation
-   **Real-time Conversion**: Instant results with smooth animations
-   **Mobile Friendly**: Responsive design works on all devices
-   **Keyboard Shortcuts**: Ctrl+Enter for quick conversions
-   **CORS Enabled**: Ready for cross-origin requests

## 📖 Usage

### Web Interface

1. **Text to Morse**: Enter English text in the first textarea and click "Convert to Morse"
2. **Morse to Text**: Enter Morse code in the second textarea and click "Convert to Text"
3. **Clear**: Use the clear buttons to reset the textareas
4. **Keyboard Shortcut**: Press `Ctrl+Enter` in any textarea to convert

### Morse Code Format

-   **Letters**: Separated by single spaces
-   **Words**: Separated by forward slashes (`/`) or double spaces
-   **Example**: `.... . .-.. .-.. --- / .-- --- .-. .-.. -.."` = "HELLO WORLD"

## 🔧 API Reference

### Convert Text/Morse

**Endpoint:** `POST /convert`

**Request Body:**

```json
{
    "text": "your text or morse code here",
    "type": "to_morse" // or "to_text"
}
```

**Response:**

```json
{
    "result": "converted text",
    "error": "error message if any"
}
```

### Examples

**Convert to Morse:**

```bash
curl -X POST http://localhost:8080/convert \
  -H "Content-Type: application/json" \
  -d '{"text": "HELLO WORLD", "type": "to_morse"}'
```

**Convert to Text:**

```bash
curl -X POST http://localhost:8080/convert \
  -H "Content-Type: application/json" \
  -d '{"text": ".... . .-.. .-.. --- / .-- --- .-. .-.. -..", "type": "to_text"}'
```

## 📊 Supported Characters

| Type            | Characters                            |
| --------------- | ------------------------------------- |
| **Letters**     | A-Z (case insensitive)                |
| **Numbers**     | 0-9                                   |
| **Punctuation** | `. , ? ' ! / ( ) & : ; = + - _ " $ @` |
| **Special**     | Space (converted to `/` in Morse)     |

## 🏗️ Project Structure

```
morse-converter-vercel/
├── api/
│   └── handler.go          # Serverless function
├── public/
│   └── index.html          # Frontend
├── vercel.json             # Vercel configuration
├── go.mod                  # Go module file
└── README.md               # Documentation
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

## 🎯 Roadmap

-   [ ] Audio playback of Morse code
-   [ ] Morse code input via sound/microphone
-   [ ] Custom timing settings
-   [ ] Batch file conversion
-   [ ] API rate limiting
-   [ ] User authentication for advanced features
-   [ ] Mobile app version

---

**Made with ❤️ and Go**

If you find this project useful, please consider giving it a ⭐ on GitHub!
