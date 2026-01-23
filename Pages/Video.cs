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
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Hosting;
using Microsoft.AspNetCore.Http;
using System.IO;
using System.ComponentModel.DataAnnotations;

namespace aspnet_webapp.Pages
{
    [Authorize("Restricted")]
    [RequestSizeLimit(long.MaxValue)]        // No request limit
    [RequestFormLimits(MultipartBodyLengthLimit = long.MaxValue)]
    
    public class VideoModel : PageModel
    {
        private readonly ILogger _logger;
        private readonly IConfiguration _configuration;
        private readonly IVideoService _videoService;


        public VideoModel(IVideoService videoService,IConfiguration configuration,ILogger<VideoModel> logger)
        {
            _logger = logger;
            _configuration=configuration;
            _videoService=videoService;
            Videos=videoService.GetVideos().ToArray();
        }

        public VideoInfo[] Videos { get; set; }
        
        [BindProperty]
        public VideoInfo SelectedVideo { get; set; }
        
        [BindProperty]
        public string SelectedVideoId { get; set; } = "";

        public void OnPostPlay() {
            SelectedVideo=Videos.FirstOrDefault(v=>v.Id==(SelectedVideoId ?? ""));
        }

        public string ConvertToMB(long value) {
            return string.Format("{0} MB", value/1024/1024);
        }


        public string UploadError { get; set; }

        public async Task OnPostUploadAsync() {
             // Access the file directly from the request collection
            var file = HttpContext.Request.Form.Files["file"];

            if (file != null && file.Length > 0)
            {
                try
                {
                    _logger.LogInformation("starting upload");
                    await Task.CompletedTask;
                    await _videoService.UploadVideoAsync(file.OpenReadStream(), file.FileName);    
                    _logger.LogInformation("finished upload");
                }
                catch (System.Exception ex)
                {
                    _logger.LogError(ex,ex.Message);
                    UploadError=ex.Message;
                }
            } else
            {
                ViewData["Message"] = "No file selected or file is empty.";
            }   
        }
    }
}
