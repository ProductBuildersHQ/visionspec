/**
 * ProductBuildersHQ Navigation Initializer
 * Version: 1.0.0
 *
 * Usage:
 *   <div id="pbhq-navbar-container"></div>
 *   <script src="https://productbuildershq.com/js/site-nav.min.js"></script>
 *   <script src="https://productbuildershq.com/js/pbhq-nav.js"></script>
 */

(function() {
  'use strict';

  // Theme CSS for ProductBuildersHQ colors - overrides design system variables
  var themeCSS = `
    :root,
    :host,
    [theme="dark"],
    :host([theme="dark"]) {
      /* ProductBuildersHQ Design System Overrides */
      /* Background - Forest green theme */
      --ds-bg-primary: #064e3b;
      --ds-bg-secondary: #0f5c4a;
      --ds-bg-tertiary: #157a5a;
      --ds-bg-inverse: #f0fdf4;

      /* Text - White on dark green */
      --ds-text-primary: #f1f5f9;
      --ds-text-secondary: #cbd5e1;
      --ds-text-muted: rgba(255, 255, 255, 0.6);
      --ds-text-inverse: #064e3b;

      /* Border */
      --ds-border-default: rgba(255, 255, 255, 0.15);
      --ds-border-strong: rgba(255, 255, 255, 0.25);
      --ds-border-subtle: rgba(255, 255, 255, 0.08);

      /* Accent - Emerald/Gold */
      --ds-accent: #10b981;
      --ds-accent-hover: #34d399;
      --ds-accent-muted: #065f46;

      /* Status */
      --ds-status-success: #34d399;
      --ds-status-warning: #facc15;
      --ds-status-error: #f87171;
      --ds-status-info: #60a5fa;
    }
  `;

  // GitHub icon
  var githubIcon = '<svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>';

  // Base URL
  var BASE_URL = 'https://productbuildershq.com';
  var GITHUB_URL = 'https://github.com/ProductBuildersHQ';

  // Navigation configuration
  var navbarConfig = {
    brand: {
      name: 'ProductBuildersHQ',
      href: BASE_URL
    },
    baseUrl: BASE_URL,
    links: [
      { id: 'frameworks', label: 'Frameworks', href: '/frameworks' },
      { id: 'case-studies', label: 'Case Studies', href: '/case-studies' }
    ],
    dropdowns: [
      {
        id: 'products',
        label: 'Products',
        items: [
          {
            id: 'visionspec',
            label: 'VisionSpec',
            href: '/visionspec/',
            description: 'Multi-domain specification orchestration'
          }
        ]
      },
      {
        id: 'resources',
        label: 'Resources',
        items: [
          { id: 'papers', label: 'Papers', href: '/papers/' },
          { id: 'about', label: 'About', href: '/about' }
        ]
      }
    ],
    actions: [
      {
        id: 'github',
        label: 'GitHub',
        href: GITHUB_URL,
        icon: githubIcon,
        external: true
      }
    ]
  };

  function init() {
    var containerId = 'pbhq-navbar-container';
    var container = document.getElementById(containerId);

    if (!container) {
      console.warn('ProductBuildersHQ Nav: Container #' + containerId + ' not found');
      return;
    }

    // Check if wt-navbar is defined (site-nav.min.js loaded)
    if (typeof customElements.get('wt-navbar') === 'undefined') {
      console.warn('ProductBuildersHQ Nav: wt-navbar not found. Make sure site-nav.min.js is loaded.');
      return;
    }

    // Inject theme CSS
    if (!document.getElementById('pbhq-theme-css')) {
      var style = document.createElement('style');
      style.id = 'pbhq-theme-css';
      style.textContent = themeCSS;
      document.head.appendChild(style);
    }

    // Create navbar
    var navbar = document.createElement('wt-navbar');
    navbar.setAttribute('theme', 'dark');
    navbar.config = navbarConfig;

    container.appendChild(navbar);
  }

  // Auto-initialize
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }

  // Export for manual initialization
  window.PbhqNav = {
    init: init,
    config: navbarConfig
  };
})();
