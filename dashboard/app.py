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

def get_data(file_path: str) -> pd.DataFrame:
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
    except Exception as e:
        logger.error(f"Une erreur est survenue lors du chargement du fichier '{file_path}': {e}")
        raise e

def main() -> None:
    """
    Fonction principale du dashboard.
    """
    st.title("Dashboard")

    # Chargement des données
    file_path = "data.csv"
    data = get_data(file_path)

    # Création des graphiques
    fig = px.line(data, x="date", y="value")
    st.plotly_chart(fig, use_container_width=True)

if __name__ == "__main__":
    main()