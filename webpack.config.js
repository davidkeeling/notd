var path = require('path');
var ExtractTextPlugin = require('extract-text-webpack-plugin');
var webpack = require('webpack');

module.exports = {
  entry: {
    main: ['./static/scripts/home.js']
  },
  output: {
    path: __dirname + '/static/scripts/build',
    filename: '[name].dev.js'
  },
  plugins: [
    new ExtractTextPlugin('styles.[name].dev.css')
  ],
  resolve: {
    modules: [
      path.resolve('./static/scripts'),
      path.resolve('./node_modules')
    ]
  },
  devtool: 'inline-source-map',
  module: {
    loaders: [
      {
        test: /\.js?$/,
        loader: 'babel-loader',
        exclude: /node_modules/
      },
      // {
      //   test: /\.css$/,
      //   loader: ExtractTextPlugin.extract('css-loader?modules&localIdentName=[name]---[local]---[hash:base64:5]'),
      //   exclude: /node_modules/
      // },
      // {
      //   //don't assign new classNames to source package css modules:
      //   test: /.*node_modules.*\.css$/,
      //   loader: ExtractTextPlugin.extract('css-loader')
      // },
      // {
      //   test: /\.scss$/,
      //   loader: ExtractTextPlugin.extract('css-loader?modules&importLoaders=1&localIdentName=[name]__[local]---[hash:base64:5]!sass-loader?{"includePaths":["./app/static/scripts/browser"]}'),
      // }
    ]
  }
};