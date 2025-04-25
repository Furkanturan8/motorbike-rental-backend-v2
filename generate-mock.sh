#!/bin/bash

# Kullanım: ./generate-mock.sh MODEL_NAME COUNT

MODEL_NAME=$1
COUNT=${2:-1}

if [[ -z "$MODEL_NAME" ]]; then
  echo "Kullanım: $0 ModelAdi [Adet]"
  echo "Örnek: $0 User 5"
  exit 1
fi

MODEL_DIR="internal/model"
MOCK_DIR="mock_data"
MOCK_FILE="${MOCK_DIR}/${MODEL_NAME}.json"

# MODEL_NAME struct'ını içeren dosyayı bul
MODEL_FILE=$(grep -rl "type ${MODEL_NAME} struct" "$MODEL_DIR")

if [[ -z "$MODEL_FILE" ]]; then
  echo "Struct bulunamadı: $MODEL_NAME"
  exit 1
fi

mkdir -p "$MOCK_DIR"

# Sadece MODEL_NAME'e ait struct'ın satırlarını al
STRUCT_LINES=$(awk "/type ${MODEL_NAME} struct/,/^\}/" "$MODEL_FILE")

echo "[" > "$MOCK_FILE"

for ((i = 0; i < COUNT; i++)); do
  echo "  {" >> "$MOCK_FILE"
  FIELD_INDEX=0

  while IFS= read -r LINE; do
    # Alan satırını ayıkla (ilişkisel alanları ve yorum satırlarını dışla)
    FIELD_LINE=$(echo "$LINE" | grep -E '^[[:space:]]+[A-Z][A-Za-z0-9_]+[[:space:]]+[^\[]+[[:space:]]+`' | grep -v 'rel:')

    if [[ -n "$FIELD_LINE" ]]; then
      FIELD_NAME=$(echo "$FIELD_LINE" | awk '{print $1}')
      FIELD_TYPE=$(echo "$FIELD_LINE" | awk '{print $2}' | sed 's/^\*//;s/\[\]//;s/^.*\.//')
      JSON_TAG=$(echo "$LINE" | grep -o 'json:"[^"]*"' | cut -d'"' -f2)
      JSON_KEY=${JSON_TAG:-$FIELD_NAME}

     # Eğer json tag'inde "-" varsa bu alanı atla
      if [[ "$JSON_TAG" == "-" ]]; then
        continue
      fi

      # Değer oluştur
      case "$FIELD_TYPE" in
        string)
          if [[ "$FIELD_NAME" == "Email" || "$JSON_KEY" == "email" ]]; then
            VALUE="\"user${i}${RANDOM}@example.com\""
          else
            VALUE="\"${JSON_KEY}_${RANDOM}\""
          fi
          ;;
        int|int64|int32)
          VALUE=$((RANDOM % 10000))
          ;;
        bool)
          VALUE=$(if [[ $((RANDOM % 2)) -eq 0 ]]; then echo "true"; else echo "false"; fi)
          ;;
        float|float64)
          VALUE=$(awk 'BEGIN{srand(); printf "%.2f", rand()*100}')
          ;;
        Time|time.Time)
          VALUE="\"$(date -Iseconds)\""
          ;;
        Role) # When you wanna add a new role, add it to the list
          VALUE="\"$(if [[ $((RANDOM % 2)) -eq 0 ]]; then echo "user"; else echo "admin"; fi)\""
          ;;
        Status) # When you wanna add a new status, add it to the list
          STATUSES=("active" "inactive" "banned")
          RANDOM_STATUS=${STATUSES[$RANDOM % ${#STATUSES[@]}]}
          VALUE="\"$RANDOM_STATUS\""
          ;;
        *)
          VALUE="null"
          ;;
      esac

      if [[ $FIELD_INDEX -gt 0 ]]; then
        echo "," >> "$MOCK_FILE"
      fi
      echo -n "    \"${JSON_KEY}\": ${VALUE}" >> "$MOCK_FILE"
      ((FIELD_INDEX++))
    fi
  done <<< "$STRUCT_LINES"

  if [[ $i -lt $((COUNT - 1)) ]]; then
    echo -e "\n  }," >> "$MOCK_FILE"
  else
    echo -e "\n  }" >> "$MOCK_FILE"
  fi
done

echo "]" >> "$MOCK_FILE"

echo "$MODEL_NAME için ${COUNT} adet mock veri oluşturuldu: $MOCK_FILE"