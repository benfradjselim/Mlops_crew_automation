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

        Raises:
        ValueError: If model_path or data_path is None.
        """
        if model_path is None or data_path is None:
            raise ValueError("Model path and data path cannot be None")
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
        max_depth (int): The maximum depth of the trees.
        min_samples_split (int): The minimum number of samples required to split an internal node.
        min_samples_leaf (int): The minimum number of samples required to be at a leaf node.

        Returns:
        RandomForestClassifier: The trained model.

        Raises:
        ValueError: If the data is empty.
        """
        if data.shape[0] == 0:
            raise ValueError("Data is empty")
        X, y = data[:, :-1], data[:, -1]
        X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)
        model = RandomForestClassifier(n_estimators=n_estimators, max_depth=max_depth, min_samples_split=min_samples_split, min_samples_leaf=min_samples_leaf)
        model.fit(X_train, y_train)
        return model

    def train(self, n_estimators: int = 100, max_depth: int = 10, min_samples_split: int = 2, min_samples_leaf: int = 1) -> None:
        """
        Train the model using the data from the data_path.

        Args:
        n_estimators (int): The number of estimators in the random forest.
        max_depth (int): The maximum depth of the trees.
        min_samples_split (int): The minimum number of samples required to split an internal node.
        min_samples_leaf (int): The minimum number of samples required to be at a leaf node.

        Raises:
        FileNotFoundError: If the data file is not found.
        ValueError: If the data is empty.
        """
        try:
            data = self._load_data()
            model = self._train_model(data, n_estimators, max_depth, min_samples_split, min_samples_leaf)
            self.model = model
            save_model(self.model_path, model)
            logger.info("Model trained and saved")
        except Exception as e:
            logger.error(f"Error training model: {e}")
            raise

    def predict(self, data: np.ndarray) -> np.ndarray:
        """
        Make predictions using the trained model.

        Args:
        data (np.ndarray): The data to make predictions on.

        Returns:
        np.ndarray: The predicted labels.

        Raises:
        ValueError: If the model is not trained.
        """
        if self.model is None:
            raise ValueError("Model is not trained")
        return self.model.predict(data)

    def evaluate(self, data: np.ndarray) -> float:
        """
        Evaluate the model using the data.

        Args:
        data (np.ndarray): The data to evaluate the model on.

        Returns:
        float: The accuracy of the model.

        Raises:
        ValueError: If the model is not trained.
        """
        if self.model is None:
            raise ValueError("Model is not trained")
        return accuracy_score(data[:, -1], self.model.predict(data[:, :-1]))