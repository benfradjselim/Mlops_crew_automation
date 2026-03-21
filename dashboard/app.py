import streamlit as st
import matplotlib.pyplot as plt
import numpy as np
import pandas as pd

# Data generation
np.random.seed(0)
data = np.random.rand(100)
cpu_usage = np.random.randint(0, 100, 100)

# Create a Streamlit app
st.title("Dashboard")

# Add a sidebar
st.sidebar.title("Options")

# Select a time range
time_range = st.sidebar.selectbox("Select a time range", ["Last hour", "Last day", "Last week"])

# Create a line plot
fig, ax = plt.subplots()
ax.plot(data)
ax.set_title("CPU usage")
ax.set_xlabel("Time")
ax.set_ylabel("CPU usage (%)")
st.pyplot(fig)

# Create a histogram
fig, ax = plt.subplots()
ax.hist(cpu_usage, bins=10)
ax.set_title("CPU usage distribution")
ax.set_xlabel("CPU usage (%)")
ax.set_ylabel("Frequency")
st.pyplot(fig)

# Create a scatter plot for anomalies
anomalies = np.random.choice(cpu_usage, 10)
fig, ax = plt.subplots()
ax.scatter(anomalies, np.random.rand(10))
ax.set_title("Anomalies")
ax.set_xlabel("CPU usage (%)")
ax.set_ylabel("Anomaly score")
st.pyplot(fig)

# Create a bar chart for alerts
alerts = np.random.choice(cpu_usage, 10)
fig, ax = plt.subplots()
ax.bar(alerts, np.random.rand(10))
ax.set_title("Alerts")
ax.set_xlabel("CPU usage (%)")
ax.set_ylabel("Alert level")
st.pyplot(fig)

# Display a table with CPU usage data
df = pd.DataFrame(cpu_usage, columns=["CPU usage (%)"])
st.table(df)

# Display a text area with alert messages
alert_messages = ["CPU usage is high", "Anomaly detected", "Alert level is critical"]
st.write("Alert messages:")
st.write("\n".join(alert_messages))