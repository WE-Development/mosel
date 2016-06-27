/***********************************************************************************************
 * User Configuration.
 **********************************************************************************************/
/** Map relative paths to URLs. */
const map:any = {
  '@angular2-material': 'vendor/@angular2-material',
  'angular2-highcharts': 'vendor/angular2-highcharts',
  'highcharts/highstock.src': 'vendor/highcharts/highstock.src.js'
};

/** User packages configuration. */
const packages:any = {
  'angular2-highcharts': {
    /*main: `index.js`,*/
    format: 'cjs',
    defaultExtension: 'js'
  },
};

const materialPkgs:string[] = [
  'core',
  'button',
  'card',
  'checkbox',
  'card',
  'grid-list',
];

materialPkgs.forEach((pkg) => {
  packages[`@angular2-material/${pkg}`] = {
    main: `${pkg}.js`,
    format: 'cjs',
    defaultExtension: 'js'
  };
});

////////////////////////////////////////////////////////////////////////////////////////////////
/***********************************************************************************************
 * Everything underneath this line is managed by the CLI.
 **********************************************************************************************/
const barrels:string[] = [
  // Angular specific barrels.
  '@angular/core',
  '@angular/common',
  '@angular/compiler',
  '@angular/http',
  '@angular/router',
  '@angular/platform-browser',
  '@angular/platform-browser-dynamic',

  // Thirdparty barrels.
  'rxjs',

  // App specific barrels.
  'app',
  'app/shared',
  /** @cli-barrel */
];

const cliSystemConfigPackages:any = {};
barrels.forEach((barrelName:string) => {
  cliSystemConfigPackages[barrelName] = {main: 'index'};
});

/** Type declaration for ambient System. */
declare var System:any;

// Apply the CLI SystemJS configuration.
System.config({
  map: {
    '@angular': 'vendor/@angular',
    'rxjs': 'vendor/rxjs',
    'main': 'main.js'
  },
  packages: cliSystemConfigPackages
});

// Apply the user's configuration.
System.config({map, packages});
