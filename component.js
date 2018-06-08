(function (global, factory) {
    typeof exports === 'object' && typeof module !== 'undefined' ? module.exports = factory() :
        typeof define === 'function' && define.amd ? define(factory) : (global.Component = factory())
}(this, (function () {
    return {
        Component: {
           /* render: function (createElement) {

            },*/
            props: {}
        },
        install: function (Vue) {
            return (function ($) {
                Vue.mixin({
                    created: function () {
                        Vue.component('Component', $.Component);
                    }
                })
            }({
                Component: this.Component
            }));
        }
    }
})));