# word-frequency

A small command-line tool that reads text from standard input and prints the most frequent words with their counts.

## Usage

```sh
echo "The cat and the hat. The end." | go run . -n 3
go run . -n 5 < essay.txt
```

Output is one `word<TAB>count` line per word:

```
the	3
and	1
cat	1
```

| Flag | Default | Description |
|---|---|---|
| `-n` | `10` | Number of results to print. Must be greater than 0. |

To build a binary instead: `go build`, then `./word-frequency -n 3 < essay.txt`.

## Behavior

- **Words** are runs of letters (any language) and digits. Apostrophes are kept only between letters or digits: `That's` → `that's`, `'hello'` → `hello`.
- **Case** is ignored: `The` and `the` count as the same word.
- **Curly apostrophes** (`’`, `‘`) and full-width `＇` are treated as `'`, so `don’t` and `don't` match.
- **Everything else** separates words, including hyphens: `post-game` → `post`, `game`.
- **Order:** highest count first; ties are sorted alphabetically.
- If `-n` is larger than the number of distinct words, all words are printed.
- Empty input prints nothing and exits successfully.
- `-n` of 0 or less prints an error to stderr and exits with status 1.

## Testing

```sh
go test ./...
```
