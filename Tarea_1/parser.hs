import Text.Parsec
import Text.Parsec.String (Parser)
import Text.Parsec.Char (char, string, digit, noneOf, oneOf, hexDigit)
import Text.Parsec.Combinator (many1, sepBy, sepBy1, optionMaybe)
import Data.Functor.Identity (Identity)

data JSONValue = JSONObject [(String, JSONValue)]
               | JSONArray [JSONValue]
               | JSONString String
               | JSONNumber Double
               | JSONBool Bool
               | JSONNull
               deriving (Show, Eq)

jsonParser :: Parser JSONValue
jsonParser = jsonObject <|> jsonArray

jsonObject :: Parser JSONValue
jsonObject = do
    char '{'
    spaces
    members <- jsonMembers
    spaces
    char '}'
    return $ JSONObject members

jsonMembers :: Parser [(String, JSONValue)]
jsonMembers = jsonPair `sepBy` (spaces >> char ',' >> spaces)

jsonPair :: Parser (String, JSONValue)
jsonPair = do
    JSONString key <- jsonString
    spaces
    char ':'
    spaces
    value <- jsonValue
    return (key, value)

jsonArray :: Parser JSONValue
jsonArray = do
    char '['
    spaces
    elements <- jsonElements
    spaces
    char ']'
    return $ JSONArray elements

jsonElements :: Parser [JSONValue]
jsonElements = jsonValue `sepBy` (spaces >> char ',' >> spaces)

jsonValue :: Parser JSONValue
jsonValue = jsonString
          <|> jsonNumber
          <|> jsonObject
          <|> jsonArray
          <|> jsonBool
          <|> jsonNull

jsonString :: Parser JSONValue
jsonString = do
    char '"'
    str <- many (noneOf "\"\\" <|> jsonEscape)
    char '"'
    return $ JSONString str

jsonEscape :: Parser Char
jsonEscape = do
    char '\\'
    esc <- oneOf "\\\"/bfnrtu"
    case esc of
        'u' -> do
            hex1 <- hexDigit
            hex2 <- hexDigit
            hex3 <- hexDigit
            hex4 <- hexDigit
            return $ read ("\\x" ++ [hex1, hex2, hex3, hex4])
        _ -> return $ case esc of
                        '\\' -> '\\'
                        '"' -> '"'
                        '/' -> '/'
                        'b' -> '\b'
                        'f' -> '\f'
                        'n' -> '\n'
                        'r' -> '\r'
                        't' -> '\t'

jsonNumber :: Parser JSONValue
jsonNumber = do
    num <- many1 (digit <|> oneOf "-+eE.")
    return $ JSONNumber (read num)

jsonBool :: Parser JSONValue
jsonBool = (string "true" >> return (JSONBool True)) <|> (string "false" >> return (JSONBool False))

jsonNull :: Parser JSONValue
jsonNull = string "null" >> return JSONNull

parseJSON :: String -> Either ParseError JSONValue
parseJSON = parse jsonParser ""

main :: IO ()
main = do
    let input = "{\"nombre\": \"Jesus\", \"usbid\": 1910072, \"estudiante\": true}"
    let result = parseJSON input
    print result

