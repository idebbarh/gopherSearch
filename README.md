# GopherSearch

GopherSearch is a search engine written in GoLang that indexes HTML files within a specified folder, creating a JSON index file. The program also serves a local server that allows users to perform searches using TF-IDF (Term Frequency-Inverse Document Frequency) ranking on the indexed files and returns relevant search results to the user.

## Installation

To install GopherSearch, make sure you have GoLang installed on your machine. You can then clone the repository:

```bash
    git clone https://github.com/yourusername/gophersearch.git
```

## Usage

To index the HTML files and serve the local server, use the following command:

```bash
    gophersearch serve <path_to_folder>
```

Replace <path_to_folder> with the path to the folder containing the HTML files you want to index.

## Searching Process

Once the server is running, you can perform searches using a web browser or API requests. The search query will be processed using TF-IDF on the indexed files. The server will return relevant search results to the user.

## Contributing

Contributions to GopherSearch are welcome! Feel free to fork the repository and submit pull requests for any improvements or new features.

## License

This project is open-source and provided under the MIT license. Please refer to the LICENSE file for detailed information regarding the terms and conditions of use, modification, and distribution.
