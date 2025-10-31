import joblib
import numpy as np
from scipy.sparse import hstack

model = joblib.load("model/ticket_classifier_final.pkl")
vecs = joblib.load("model/vectorizers_final.pkl")

subject_vectorizer = vecs["subject_vectorizer"]
body_vectorizer = vecs["body_vectorizer"]
scaler = vecs["scaler"]

def predict_category(subject: str, body: str):
   
    X_sub = subject_vectorizer.transform([subject])
    X_body = body_vectorizer.transform([body])


    len_subject = len(subject.split())
    len_body = len(body.split())
    has_urgent = int("urgent" in body.lower())
    has_refund = int("refund" in body.lower())
    has_password = int("password" in body.lower())
    numeric_features = np.array([[len_subject, len_body, has_urgent, has_refund, has_password]])

    X_num_scaled = scaler.transform(numeric_features)

   
    X_final = hstack([X_sub, X_body, X_num_scaled])

   
    pred = model.predict(X_final)[0]
    return pred
