const path = require('path');
const { VueLoaderPlugin } = require('vue-loader')
const autoprefixer = require('autoprefixer');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");
const TerserPlugin = require("terser-webpack-plugin");

const isDevServer = process.env.WEBPACK_DEV_SERVER || process.env.WEBPACK_SERVE;

module.exports = {
  entry: './src/main.js',
  output: {
    path: path.resolve(__dirname, 'ui'),
    filename: 'bundle.js',
    publicPath: isDevServer ? "auto" : "/ui/",
  },
  module: {
    rules: [
      {
        test: /\.css$/i,
        use: [
          {
            loader: MiniCssExtractPlugin.loader,
          },
          {
            loader: 'css-loader',
            options: {
              sourceMap: true
            }
          },
          {
            loader: 'postcss-loader',
            options: {
              sourceMap: true,
              postcssOptions: {
                plugins: () => [autoprefixer()],
              }
            }
          },

          // {
          //   loader: 'sass-loader',
          //   options: {
          //     sourceMap: true,
          //     sassOptions: {
          //       includePaths: [
          //         path.resolve(__dirname, "node_modules")
          //       ]
          //     }
          //   }
          // },
        ]
      },
      {
        test: /\.vue$/,
        loader: 'vue-loader'
      },
      {
        test: /\.(mov|mp4|mp3|wav|ogg|pdf)$/i,
        type: "asset/resource",
        generator: {
          filename: 'assets/videos/[hash][ext][query]'
        },
      },
      {
        test: /\.(eot|svg|ttf|woff|woff2)$/i,
        type: 'asset/resource',
        generator: {
          filename: 'assets/fonts/[hash][ext][query]'
        },
        parser: {
          dataUrlCondition: {
            maxSize: 4 * 1024 // 4kb
          }
        },
      },
      {
        test: /\.(|png|jpe?g|gif)$/i,
        type: 'asset',
        generator: {
          filename: 'assets/img/[hash][ext][query]'
        },
      },
    ]
  },
  plugins: [
    new MiniCssExtractPlugin({
      // Options similar to the same options in webpackOptions.output
      // all options are optional
      filename: 'css/[name].css',
      // chunkFilename: 'css/[id].css',
      ignoreOrder: false, // Enable to remove warnings about conflicting order
    }),
    new VueLoaderPlugin(),
    new HtmlWebpackPlugin({
      template: path.resolve(__dirname, 'src/index.html'),
    })
  ],
  devServer: {
    open: true,
    compress: true,
    port: 9000,
    allowedHosts: [
      "localhost",
      ".demo2.mixmedia.com",
    ],
    proxy: {
      "/api": {
        target: "http://127.0.0.1:8843",
        // secure: true,
        changeOrigin: true,
      },
    },
  },
  devtool: "source-map",
  watchOptions: {
    ignored: /node_modules/
  },
  optimization: {
    minimizer: [
      new CssMinimizerPlugin({
        minimizerOptions: {
          preset: [
            'default',
            {
              mergeLonghand: false,
              cssDeclarationSorter: false
            }
          ]
        },
      }),
      new TerserPlugin(),
    ],
    minimize: process.env.NODE_ENV !== 'development',
  },
  mode: process.env.NODE_ENV,
};
