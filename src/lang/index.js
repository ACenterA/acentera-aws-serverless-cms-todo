import Vue from 'vue'
import VueI18n from 'vue-i18n'
import Cookies from 'js-cookie'
import elementEnLocale from 'element-ui/lib/locale/lang/en' // element-ui lang
import elementFrLocale from 'element-ui/lib/locale/lang/fr' // element-ui lang
import enLocale from './en'
import frLocale from './fr'

Vue.use(VueI18n)

const messages = {
  en: {
    ...enLocale,
    ...elementEnLocale
  },
  fr: {
    ...frLocale,
    ...elementFrLocale
  }
}

const i18n = new VueI18n({
  // set locale
  // options: en or zh
  locale: Cookies.get('language') || 'en',
  // set locale messages
  messages
})

/* Extend i18n messages */
if (window.app._i18n) {
  for (var k in messages) {
    window.app._i18n.mergeLocaleMessage(k, messages[k])
  }
}

export default i18n
