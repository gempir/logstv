const path = require("path");

const SRC_DIR = path.resolve(__dirname, 'src');
const BUILD_DIR = path.resolve(__dirname, 'build');

module.exports = {
    entry: './src/index.jsx',
    output: {
		path: BUILD_DIR,
		filename: 'bundle.js',
		publicPath: "/",
	},
	module: {
		rules: [
		    {
				test: /\.jsx$/,
				exclude: /node_modules/,
				use: {
					loader: "babel-loader"
				}
			},
			{
				test: /\.s?css$/,
				use: [{
					loader: "style-loader"
				}, {
					loader: "css-loader", options: {
						sourceMap: true
					}
				}, {
					loader: "sass-loader", options: {
						sourceMap: true
					}
				}]
			}
		]
	},
	resolve: {
		extensions: ['.js', '.jsx'],
	}
};