# compare algorithms
import pandas as pd
from pandas import read_csv
from matplotlib import pyplot
from sklearn.ensemble import RandomForestClassifier, ExtraTreesClassifier
from sklearn.model_selection import train_test_split
from sklearn.model_selection import cross_val_score
from sklearn.model_selection import StratifiedKFold
from sklearn.linear_model import LogisticRegression
from sklearn.neural_network import MLPClassifier
from sklearn.preprocessing import LabelEncoder
from sklearn.tree import DecisionTreeClassifier, ExtraTreeClassifier
from sklearn.neighbors import KNeighborsClassifier, RadiusNeighborsClassifier
from sklearn.discriminant_analysis import LinearDiscriminantAnalysis
from sklearn.naive_bayes import GaussianNB, BernoulliNB
from sklearn.svm import LinearSVC

def encode_column(target_feature):
    encoder = LabelEncoder()
    cols = dataset[target_feature].values
    encoded_values = encoder.fit_transform(cols)
    dataset[target_feature] = pd.Series(encoded_values, index=dataset.index)

# Load dataset
dataset = read_csv("../../data/historical_data_cleaned.csv")
dataset['Country1RankConfederation'] = dataset['Country1RankConfederation'].fillna('UNKNOWN')
dataset['Country2RankConfederation'] = dataset['Country1RankConfederation'].fillna('UNKNOWN')
dataset = dataset.fillna(0)

encode_column('Country1')
encode_column('Country2')
encode_column('Country1RankConfederation')
encode_column('Country2RankConfederation')
encode_column('MatchDate')
encode_column('Tournament')
encode_column('Neutral')


# Split-out validation dataset
X = dataset.drop(['Country1Score','Country2Score'], 'columns')
Y = dataset[['Country1Score','Country2Score']]

X = X.values
y = Y.values

X_train, X_validation, Y_train, Y_validation = train_test_split(X, y, test_size=0.05, shuffle=False)

# Spot Check Algorithms
models = []
# models.append(('KNeighborsClassifier', KNeighborsClassifier(algorithm='auto')))
# models.append(('DecisionTreeClassifier', DecisionTreeClassifier()))
models.append(('RandomForestClassifier gini', RandomForestClassifier(criterion="gini")))
# models.append(('RandomForestClassifier entropy', RandomForestClassifier(criterion="entropy")))
# models.append(('ExtraTreesClassifier', ExtraTreesClassifier()))
# models.append(['ExtraTreeClassifier', ExtraTreeClassifier()])

# evaluate each model in turn
results = []
names = []
for name, model in models:
    model.fit(X_train, Y_train)
    prediction = model.predict(X_validation)
    results.append((name, prediction))


for name, prediction in results:
    c1_correct_count = 0
    c2_correct_count = 0
    both_correct_count = 0
    for idx, result in enumerate(Y_validation):
        c1_result = result[0]
        c2_result = result[1]
        score = prediction[idx]
        c1_score = score[0]
        c2_score = score[1]
        if c1_result == c1_score:
            c1_correct_count += 1
        if c2_result == c2_score:
            c2_correct_count += 1
        if c1_result == c1_score and c2_result == c2_score:
            both_correct_count += 1

    c1_accuracy = c1_correct_count/len(Y_validation)
    c2_accuracy = c2_correct_count/len(Y_validation)
    accuracy = both_correct_count/len(Y_validation)
    print("\n{} -> Country 1 Accuracy: {}".format(name, c1_accuracy))
    print("{} -> Country 2 Accuracy: {}".format(name, c2_accuracy))
    print("{} -> Both Scores Accuracy: {}\n".format(name, accuracy))







#%%
