using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using aspnet_webapp.Services;
using Azure.Security.KeyVault.Secrets;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;

namespace aspnet_webapp.Pages
{
    [Authorize("Restricted")]
    public class PrivacyModel : PageModel
    {
        private readonly ILogger<PrivacyModel> _logger;
        private readonly IUserInfoService _userService;

        public PrivacyModel(ILogger<PrivacyModel> logger, IUserInfoService userService)
        {
            _logger = logger;
            _userService = userService;
        }

        public UserInfo UserInfo { get; set; }
        public async Task OnGetAsync()
        {
            UserInfo = await _userService.GetUserInfo();
        }
    }
}
