import pandas as pd
import numpy as np
import re
from sklearn.model_selection import train_test_split, StratifiedKFold, GridSearchCV
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.linear_model import LogisticRegression
from sklearn.preprocessing import StandardScaler
from sklearn.metrics import classification_report
from sklearn.utils import resample
from scipy.sparse import hstack
import joblib

# датасет
DATA_PATH = "hf://datasets/ale-dp/bilingual-ticket-classification/data/train-00000-of-00001.parquet"



df = pd.read_parquet(DATA_PATH)
print(f"oaded {len(df)} rows total.")

# только англ
df_en = df[df["language"] == "en"].copy()

# очистка текста
def clean_text(text):
    text = text.lower()
    text = re.sub(r"http\S+|www\S+", " ", text)
    text = re.sub(r"[^a-z\s]", " ", text)
    text = re.sub(r"\s+", " ", text).strip()
    return text

df_en["subject_clean"] = df_en["subject"].fillna("").apply(clean_text)
df_en["body_clean"] = df_en["body"].fillna("").apply(clean_text)

# убираем короткие записи
df_en = df_en[(df_en["subject_clean"].str.split().str.len() > 1) |
              (df_en["body_clean"].str.split().str.len() > 3)]

#Балансировка классов (oversampling) ===

df_balanced = (
    df_en.groupby("queue", group_keys=False)
         .apply(lambda x: resample(x, replace=True, n_samples=min(len(df_en), max(200, len(x))), random_state=42))
         .reset_index(drop=True)
)
# print(f"✅ After balancing: {df_balanced['queue'].value_counts().to_dict()}")


df_balanced["len_subject"] = df_balanced["subject_clean"].str.split().str.len()
df_balanced["len_body"] = df_balanced["body_clean"].str.split().str.len()
df_balanced["has_urgent"] = df_balanced["body_clean"].str.contains("urgent").astype(int)
df_balanced["has_refund"] = df_balanced["body_clean"].str.contains("refund").astype(int)
df_balanced["has_password"] = df_balanced["body_clean"].str.contains("password").astype(int)

numeric_features = df_balanced[["len_subject", "len_body", "has_urgent", "has_refund", "has_password"]].to_numpy()


labels = df_balanced["queue"].tolist()

X_train_sub, X_test_sub, X_train_body, X_test_body, X_train_num, X_test_num, y_train, y_test = train_test_split(
    df_balanced["subject_clean"],
    df_balanced["body_clean"],
    numeric_features,
    labels,
    test_size=0.2,
    random_state=42,
    stratify=labels
)

subject_vectorizer = TfidfVectorizer(max_features=7000, ngram_range=(1,2), stop_words="english", min_df=2)
body_vectorizer = TfidfVectorizer(max_features=15000, ngram_range=(1,2), stop_words="english", min_df=2)

X_train_sub_vec = subject_vectorizer.fit_transform(X_train_sub)
X_test_sub_vec = subject_vectorizer.transform(X_test_sub)
X_train_body_vec = body_vectorizer.fit_transform(X_train_body)
X_test_body_vec = body_vectorizer.transform(X_test_body)

scaler = StandardScaler(with_mean=False)
X_train_num_scaled = scaler.fit_transform(X_train_num)
X_test_num_scaled = scaler.transform(X_test_num)


X_train_final = hstack([X_train_sub_vec, X_train_body_vec, X_train_num_scaled])
X_test_final = hstack([X_test_sub_vec, X_test_body_vec, X_test_num_scaled])

# ===  Logistic Regression ===


param_grid = {"C": [0.1, 0.5, 1.0, 2.0]}
logreg = LogisticRegression(
    max_iter=2000,
    class_weight="balanced",
    solver="saga",
    penalty="l2"
)

cv = StratifiedKFold(n_splits=5, shuffle=True, random_state=42)
grid = GridSearchCV(logreg, param_grid, cv=cv, n_jobs=-1, verbose=1)
grid.fit(X_train_final, y_train)

best_model = grid.best_estimator_
print(f"✅ Best parameters: {grid.best_params_}")


y_pred = best_model.predict(X_test_final)
# print("\n Classification Report:")
# print(classification_report(y_test, y_pred))

# сохранение
joblib.dump(best_model, "ticket_classifier_final.pkl")
joblib.dump({
    "subject_vectorizer": subject_vectorizer,
    "body_vectorizer": body_vectorizer,
    "scaler": scaler
}, "vectorizers_final.pkl")


