﻿using System;
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

namespace aspnet_webapp.Pages
{
    [Authorize("Restricted")]
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
            // var videoPath=_configuration["VIDEO_PATH"];
            // if(videoPath?.Length>0) {
            //     Videos=System.IO.Directory.GetFiles(_configuration["VIDEO_PATH"]).Select(f=>new System.IO.FileInfo(f).Name).ToArray();
            // }
        }

        public VideoInfo[] Videos { get; set; }
        
        [BindProperty]
        public VideoInfo SelectedVideo { get; set; }
        
        [BindProperty]
        public string SelectedVideoId { get; set; } = "";

        public void OnPost() {
            SelectedVideo=Videos.FirstOrDefault(v=>v.Id==(SelectedVideoId ?? ""));
        }
    }
}