var gulp = require('gulp');
var browserify = require('gulp-browserify');
var rename = require('gulp-rename');
var sourcemaps = require('gulp-sourcemaps');
var vueify = require('vueify');

gulp.task('default', function() {
  return gulp.src('src/main.js')
    .pipe(sourcemaps.init())
    .pipe(browserify({
      transform: [[{_flags: {debug: true}}, vueify]],
      standalone: 'Vue',
    }))
    .pipe(rename('js/bundle.js'))
    .pipe(sourcemaps.write('.'))
    .pipe(gulp.dest('.'));
});
