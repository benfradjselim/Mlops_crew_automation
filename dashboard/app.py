import streamlit as st
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import plotly.express as px
import logging

# Set up logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def load_data(file_path: str = 'data.csv') -> pd.DataFrame:
    """
    Load data from a CSV file.

    Args:
        file_path (str, optional): Path to the CSV file. Defaults to 'data.csv'.

    Returns:
        pd.DataFrame: Loaded data.
    """
    try:
        data = pd.read_csv(file_path)
        logger.info(f"Loaded data from {file_path}")
        return data
    except FileNotFoundError:
        logger.error(f"File not found: {file_path}")
        return None
    except pd.errors.EmptyDataError:
        logger.error(f"File is empty: {file_path}")
        return None
    except pd.errors.ParserError:
        logger.error(f"Error parsing file: {file_path}")
        return None

def load_alerts(file_path: str = 'alerts.csv') -> pd.DataFrame:
    """
    Load alerts from a CSV file.

    Args:
        file_path (str, optional): Path to the CSV file. Defaults to 'alerts.csv'.

    Returns:
        pd.DataFrame: Loaded alerts.
    """
    try:
        alerts = pd.read_csv(file_path)
        logger.info(f"Loaded alerts from {file_path}")
        return alerts
    except FileNotFoundError:
        logger.error(f"File not found: {file_path}")
        return None
    except pd.errors.EmptyDataError:
        logger.error(f"File is empty: {file_path}")
        return None
    except pd.errors.ParserError:
        logger.error(f"Error parsing file: {file_path}")
        return None

def load_graphiques(file_path: str = 'graphiques.csv') -> pd.DataFrame:
    """
    Load graphiques from a CSV file.

    Args:
        file_path (str, optional): Path to the CSV file. Defaults to 'graphiques.csv'.

    Returns:
        pd.DataFrame: Loaded graphiques.
    """
    try:
        graphiques = pd.read_csv(file_path)
        logger.info(f"Loaded graphiques from {file_path}")
        return graphiques
    except FileNotFoundError:
        logger.error(f"File not found: {file_path}")
        return None
    except pd.errors.EmptyDataError:
        logger.error(f"File is empty: {file_path}")
        return None
    except pd.errors.ParserError:
        logger.error(f"Error parsing file: {file_path}")
        return None

def main():
    """
    Main function.
    """
    # Set page title
    st.title('Dashboard')

    # Load data
    data = load_data()
    if data is None:
        st.error("Error loading data")
        return

    # Load alerts
    alerts = load_alerts()
    if alerts is None:
        st.error("Error loading alerts")
        return

    # Load graphiques
    graphiques = load_graphiques()
    if graphiques is None:
        st.error("Error loading graphiques")
        return

    # Create sidebar
    st.sidebar.title('Options')
    options = st.sidebar.selectbox('Select an option', ['Graphiques', 'Alertes', 'Data'])

    # Display graphiques
    if options == 'Graphiques':
        st.title('Graphiques')
        fig, ax = plt.subplots()
        sns.heatmap(graphiques, ax=ax)
        st.pyplot(fig)

    # Display alertes
    elif options == 'Alertes':
        st.title('Alertes')
        st.write(alerts)

    # Display data
    elif options == 'Data':
        st.title('Data')
        st.write(data)

if __name__ == '__main__':
    main()