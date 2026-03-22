import streamlit as st
import pandas as pd
import plotly.express as px
import logging

# Configuration du logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def load_data(file_path: str) -> pd.DataFrame:
    """
    Charge les données à partir du fichier CSV.

    Args:
        file_path (str): Chemin du fichier CSV.

    Returns:
        pd.DataFrame: Les données chargées.

    Raises:
        FileNotFoundError: Si le fichier n'existe pas.
        pd.errors.EmptyDataError: Si le fichier est vide.
    """
    try:
        return pd.read_csv(file_path)
    except FileNotFoundError as e:
        logger.error(f"Le fichier '{file_path}' n'existe pas.")
        raise e
    except pd.errors.EmptyDataError as e:
        logger.error(f"Le fichier '{file_path}' est vide.")
        raise e

def load_stock_data(ticker: str) -> pd.DataFrame:
    """
    Charge les données de la bourse pour le ticker spécifié.

    Args:
        ticker (str): Le ticker de la bourse.

    Returns:
        pd.DataFrame: Les données de la bourse.

    Raises:
        yf.TickerDataUnavailable: Si les données ne sont pas disponibles.
    """
    try:
        return yf.download(tickers=ticker, period='1d')
    except yf.TickerDataUnavailable as e:
        logger.error(f"Les données pour le ticker '{ticker}' ne sont pas disponibles.")
        raise e

def load_alert_data(file_path: str) -> pd.DataFrame:
    """
    Charge les données d'alerte à partir du fichier CSV.

    Args:
        file_path (str): Chemin du fichier CSV.

    Returns:
        pd.DataFrame: Les données d'alerte chargées.

    Raises:
        FileNotFoundError: Si le fichier n'existe pas.
        pd.errors.EmptyDataError: Si le fichier est vide.
    """
    try:
        return pd.read_csv(file_path)
    except FileNotFoundError as e:
        logger.error(f"Le fichier '{file_path}' n'existe pas.")
        raise e
    except pd.errors.EmptyDataError as e:
        logger.error(f"Le fichier '{file_path}' est vide.")
        raise e

def send_alert(message: str) -> None:
    """
    Envoie une alerte via l'API.

    Args:
        message (str): Le message de l'alerte.
    """
    try:
        requests.post('https://api.example.com/alert', json={'message': message})
        logger.info("L'alerte a été envoyée avec succès.")
    except requests.exceptions.RequestException as e:
        logger.error(f"Erreur lors de l'envoi de l'alerte : {e}")

def main() -> None:
    """
    Fonction principale du script.
    """
    st.title("Dashboard")

    # Chargement des données
    data = load_data('data.csv')
    stock_data = load_stock_data('AAPL')
    alert_data = load_alert_data('alert.csv')

    # Affichage des données
    st.subheader("Données")
    st.write(data)

    st.subheader("Données de la bourse")
    st.write(stock_data)

    st.subheader("Données d'alerte")
    st.write(alert_data)

    # Envoi d'une alerte
    message = "Alerte envoyée"
    send_alert(message)

if __name__ == "__main__":
    main()