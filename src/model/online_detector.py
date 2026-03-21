import numpy as np
from sklearn.ensemble import RandomForestClassifier
from sklearn.model_selection import train_test_split
from sklearn.metrics import accuracy_score
from src.model.utils import load_data, save_model, load_model
import logging

# Set up logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class OnlineDetector:
    """
    A class for online detection using a random forest classifier.

    Attributes:
    model_path (str): The path to the model file.
    data_path (str): The path to the data file.
    model (RandomForestClassifier): The trained model.
    """

    def __init__(self, model_path: str, data_path: str):
        """
        Initialize the OnlineDetector instance.

        Args:
        model_path (str): The path to the model file.
        data_path (str): The path to the data file.
        """
        self.model_path = model_path
        self.data_path = data_path
        self.model = None

    def train(self, n_estimators: int = 100, max_depth: int = 10, min_samples_split: int = 2, min_samples_leaf: int = 1):
        """
        Train the model using the data from the data_path.

        Args:
        n_estimators (int): The number of estimators in the random forest.
        max_depth (int): The maximum depth of the trees in the random forest.
        min_samples_split (int): The minimum number of samples required to split an internal node.
        min_samples_leaf (int): The minimum number of samples required to be at a leaf node.

        Returns:
        None
        """
        try:
            # Load the data
            X, y = load_data(self.data_path)
            if X is None or y is None:
                logger.error("Failed to load data")
                return

            # Split the data into training and testing sets
            X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

            # Train the model
            self.model = RandomForestClassifier(n_estimators=n_estimators, max_depth=max_depth, 
                                                 min_samples_split=min_samples_split, min_samples_leaf=min_samples_leaf)
            self.model.fit(X_train, y_train)

            # Evaluate the model
            y_pred = self.model.predict(X_test)
            accuracy = accuracy_score(y_test, y_pred)
            logger.info(f"Model accuracy: {accuracy:.3f}")

            # Save the model
            save_model(self.model, self.model_path)

        except Exception as e:
            logger.error(f"Failed to train model: {str(e)}")

    def predict(self, data: np.ndarray):
        """
        Make predictions using the trained model.

        Args:
        data (np.ndarray): The data to make predictions on.

        Returns:
        np.ndarray: The predicted labels.
        """
        try:
            if self.model is None:
                self.model = load_model(self.model_path)
            if data is None:
                logger.error("Data is None")
                return None
            return self.model.predict(data)

        except Exception as e:
            logger.error(f"Failed to make predictions: {str(e)}")
            return None

    def update(self, new_data: np.ndarray):
        """
        Update the model using the new data.

        Args:
        new_data (np.ndarray): The new data to update the model with.

        Returns:
        None
        """
        try:
            if self.model is None:
                self.model = load_model(self.model_path)
            if new_data is None:
                logger.error("New data is None")
                return
            self.model.fit(new_data, np.zeros(len(new_data)))
            save_model(self.model, self.model_path)

        except Exception as e:
            logger.error(f"Failed to update model: {str(e)}")