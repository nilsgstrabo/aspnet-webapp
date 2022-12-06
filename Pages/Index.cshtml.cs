using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Azure.Security.KeyVault.Secrets;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;

namespace aspnet_webapp.Pages
{
    public class IndexModel : PageModel
    {
        private readonly ILogger<IndexModel> _logger;
        private readonly SecretClient _secretClient;

        public IndexModel(SecretClient secretClient,ILogger<IndexModel> logger)
        {
            _logger = logger;
            _secretClient=secretClient;
        }

        public IEnumerable<string> Secrets;

        public void OnGet()
        {
            try
            {
                var secrets = _secretClient.GetPropertiesOfSecrets();
                Secrets=secrets.Select(s=>s.Name).ToList();
            }
            catch (System.Exception ex)
            {
                Secrets=new List<string>(){"something went wrong"};
                _logger.LogError(ex,ex.Message);
            }

            foreach (var item in this.Request.Headers)
            {
                _logger.LogInformation($"{item.Key}: {item.Value.ToString()}");
            }
        }
    }
}
