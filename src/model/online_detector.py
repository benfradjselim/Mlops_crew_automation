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
    A class for online detection using a HalfSpaceTrees model.

    Attributes:
    model_path (str): The path to the model file.
    data_path (str): The path to the data file.
    model (HalfSpaceTrees): The trained model.
    """

    def __init__(self, model_path: str, data_path: str) -> None:
        """
        Initialize the OnlineDetector instance.

        Args:
        model_path (str): The path to the model file.
        data_path (str): The path to the data file.

        Raises:
        ValueError: If model_path or data_path is None.
        """
        if not model_path or not data_path:
            raise ValueError("Model path and data path cannot be empty")
        self.model_path = model_path
        self.data_path = data_path
        self.model = None

    def _load_data(self) -> np.ndarray:
        """
        Load the data from the data_path.

        Returns:
        np.ndarray: The loaded data.

        Raises:
        FileNotFoundError: If the data file is not found.
        """
        try:
            data = load_data(self.data_path)
            return data
        except FileNotFoundError as e:
            logger.error(f"Data file not found: {e}")
            raise

    def _train_model(self, data: np.ndarray, n_estimators: int = 100, max_depth: int = 10, min_samples_split: int = 2, min_samples_leaf: int = 1) -> RandomForestClassifier:
        """
        Train the model using the data.

        Args:
        data (np.ndarray): The data to train the model.
        n_estimators (int): The number of estimators in the random forest.
        max_depth (int): The maximum depth of the tree.
        min_samples_split (int): The minimum number of samples required to split an internal node.
        min_samples_leaf (int): The minimum number of samples required to be at a leaf node.

        Returns:
        RandomForestClassifier: The trained model.

        Raises:
        ValueError: If n_estimators, max_depth, min_samples_split or min_samples_leaf is invalid.
        """
        if n_estimators <= 0 or max_depth <= 0 or min_samples_split <= 0 or min_samples_leaf <= 0:
            raise ValueError("Invalid hyperparameters")
        try:
            model = RandomForestClassifier(n_estimators=n_estimators, max_depth=max_depth, min_samples_split=min_samples_split, min_samples_leaf=min_samples_leaf)
            model.fit(data)
            return model
        except Exception as e:
            logger.error(f"Error training model: {e}")
            raise

    def train(self, data: np.ndarray, n_estimators: int = 100, max_depth: int = 10, min_samples_split: int = 2, min_samples_leaf: int = 1) -> None:
        """
        Train the model using the data.

        Args:
        data (np.ndarray): The data to train the model.
        n_estimators (int): The number of estimators in the random forest.
        max_depth (int): The maximum depth of the tree.
        min_samples_split (int): The minimum number of samples required to split an internal node.
        min_samples_leaf (int): The minimum number of samples required to be at a leaf node.
        """
        try:
            self.model = self._train_model(data, n_estimators, max_depth, min_samples_split, min_samples_leaf)
            save_model(self.model_path, self.model)
        except Exception as e:
            logger.error(f"Error training model: {e}")

    def predict(self, data: np.ndarray) -> np.ndarray:
        """
        Make predictions using the trained model.

        Args:
        data (np.ndarray): The data to make predictions on.

        Returns:
        np.ndarray: The predicted labels.
        """
        try:
            if not self.model:
                raise ValueError("Model not trained")
            return self.model.predict(data)
        except Exception as e:
            logger.error(f"Error making predictions: {e}")
            raise