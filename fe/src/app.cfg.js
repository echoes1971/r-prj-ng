export const app_cfg = {
    site_title: process.env.REACT_APP_SITE_TITLE > '' ? process.env.REACT_APP_SITE_TITLE : 'R-Prj'
    ,endpoint: process.env.REACT_APP_ENDPOINT > '' ? process.env.REACT_APP_ENDPOINT : "/lippolippi" // process.env.NODE_ENV !== 'production' ? "http://localhost:8080/jsonserver.php" : "https://rprj.roccoangeloni.ch/jsonserver.php"
    ,app_home_object_id: process.env.REACT_APP_HOME_OBJECT_ID > '' ? process.env.REACT_APP_HOME_OBJECT_ID : '-10'
};
